// +build cgo

/* Copyright (c) 2016 Jason Ish
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 *
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED ``AS IS'' AND ANY EXPRESS OR IMPLIED
 * WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT,
 * INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT,
 * STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING
 * IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */

package sqlite

import (
	"github.com/jasonish/evebox/eve"
	"github.com/jasonish/evebox/log"
	"time"
)

type SqlitePurger struct {
	db     *SqliteService
	period string
	limit int64
}

func (p *SqlitePurger) Run() {
	if p.period == "0" {
		return
	}
	for {
		count, _ := p.Purge()
		if count < p.limit {
			time.Sleep(1 * time.Minute)
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (p *SqlitePurger) Purge() (int64, error) {
	period, err := time.ParseDuration(p.period)
	if err != nil {
		log.Error("%v", err)
		return 0, err
	}

	now := time.Now()
	then := eve.FormatTimestamp(now.Add(- period))
	log.Info("Deleting events prior to %v", then)

	tx, err := p.db.GetTx()
	if err != nil {
		log.Error("%v", err)
		return 0, err
	}
	defer tx.Rollback()

	start := time.Now()

	// Wrapping in a subselect so we can limit the number of events
	// deleted per run.
	q := `
delete
from events
where rowid in
    (select rowid
     from events
     where timestamp < ?
     and escalated = 0
     limit ?)`
	r, err := tx.Exec(q, then, p.limit)
	if err != nil {
		log.Error("%v", err)
		return 0, err
	}

	count, err := r.RowsAffected()
	if err != nil {
		log.Warning("Failed to get number of events purged")
	}

	err = tx.Commit()
	if err != nil {
		log.Error("%v", err)
		return 0, err
	}

	log.Info("Purged %d events in %v", count, time.Now().Sub(start))
	return count, nil
}
