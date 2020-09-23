package regip

import (
	"errors"
)

var (
	AlreadyIndexed = errors.New("already indexed")
	CastFailed     = errors.New("cast from resource to record failed")
)

const (
	INDEX_BATCH_SIZE = 100
)

func (db *DB) IndexRecords(ids []ID, lg *Logger) {
	lgr := lg.Tag("ir", CLR_indexrecords)
	lgr.Printf("start with %d records", len(ids))
	defer lgr.Print("end")

	tris := make(map[string][]ID)
	for _, cur := range ids {
		// Save the type byte
		ct := cur[0]
		// Change to index record to check existence...
		cur.SetType(RT_indexRecord)
		// Check if it's already indexed
		ex := db.Exists(cur)
		// Don't forget to change back...
		cur.SetType(ct)
		if ex {
			lgr.Printf("IndexRecord for %s exists, skipping...", cur)
			continue // Nothing to do here...
		}

		res, err := db.Get(cur, lgr)
		if err != nil {
			// FIXME: LOG
			lgr.Errorf("Trying to get %s got error %s, skipping...", cur, err)
			continue
		}

		rec, ok := res.(*Record)
		if !ok {
			lgr.Errorf("Casting %s to record failed, skipping...", cur)
			continue
		}
		rid := rec.ID()
		combis := GenerateTrigrams(rec.Name, true)

		for _, cmb := range combis {
			tris[cmb] = append(tris[cmb], rid)
		}

		// Add index records
		ir := NewIndexRecord(cur)
		err = db.Add(ir)
		if err != nil {
			lgr.Error("Failed to mark index record for ", cur)
		} else {
			lgr.Print("Added index record ", ir)
		}
	}
	// Add tris to trigrams
	for cmb, ds := range tris {
		curTri := NewTrigram(cmb, ds)
		curID := curTri.ID()
		lgr.Printf("Adding trigram %s with %d new ids", curID.String(), len(curTri.Ids))
		if db.Exists(curID) {
			old, err := db.Get(curID, lgr)
			if err != nil {
				// If this happens -- something's wrong
				panic(err)
			}
			old.(*Trigram).AddIds(curTri.Ids)
			db.Update(old)
		} else {
			db.Add(curTri)
		}
	}
}

func (db *DB) FullIndex(lg *Logger) {
	lgr := lg.Tag("fullindex", CLR_fullindex)

	idf := db.IDFlow(RT_record, lgr)
	var buf []ID
	for {
		cur, ok := idf.Get()
		if !ok {
			lgr.Print("IDFlow.Get was not ok, exiting...")
			break
		}
		buf = append(buf, cur.(ID))
		if len(buf) >= INDEX_BATCH_SIZE {
			lgr.Printf("calling IndexRecords with %d records", len(buf))
			db.IndexRecords(buf, lgr)
			buf = []ID{}
		}
	}
	if len(buf) > 0 {
		lgr.Printf("Finishing up index with %d records", len(buf))
		db.IndexRecords(buf, lgr)
	}
}
