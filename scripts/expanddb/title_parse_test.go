package main

import "testing"

func TestTitleParsing(t *testing.T) {
	cases := []struct {
		Input          ArchiveEntry
		ExpectedResult ArchiveEntry
	}{
		// {
		// 	Input: ArchiveEntry{
		// 		Title: "mogucomi  第SP回 ゆみりと愛奈のモグモグ・コミュニケーションズ",
		// 	},
		// 	ExpectedResult: ArchiveEntry{
		// 		EpTitle: "第SP回",
		// 		Title:   "ゆみりと愛奈のモグモグ・コミュニケーションズ",
		// 	},
		// },
		// {
		// 	Input: ArchiveEntry{
		// 		Title: "stbsp  第SP回 『ストライク・ザ・ブラッド恩萊祭 -ONLINE FES- 』開催直前スペシャルラジオ",
		// 	},
		// 	ExpectedResult: ArchiveEntry{
		// 		EpTitle: "第SP回",
		// 		Title:   "『ストライク・ザ・ブラッド恩萊祭 -ONLINE FES- 』開催直前スペシャルラジオ",
		// 	},
		// },
		{
			Input: ArchiveEntry{Size: 31786337, MessageID: 3141, Duration: 1957, Title: "200707 gochiusabloom ご注文はラジオですか？BLOOM", Performer: "ココア 役 佐倉綾音 チノ 役 水瀬いのり", Filename: "200707-gochiusabloom--onsen.mp3", Date: "200707", ShowID: "gochiusabloom", Source: "", EpTitle: "", PartIdx: 0},
			ExpectedResult: ArchiveEntry{
				EpTitle: "",
				Title:   "ご注文はラジオですか？BLOOM",
			},
		},
		{
			Input: ArchiveEntry{Size: 19648588, MessageID: 3150, Duration: 1200, Title: "190322 gridman p1 アニメGRIDMAN ラジオ とりあえずUNION", Performer: "響裕太 役 広瀬裕也 宝多六花 役 宮本侑芽", Filename: "190322-gridman--onsen-p1.mp3", Date: "190322", ShowID: "gridman", Source: "", EpTitle: "", PartIdx: 0},
			ExpectedResult: ArchiveEntry{
				Title:   "アニメGRIDMAN ラジオ とりあえずUNION",
				PartIdx: 1,
			},
		},
	}
	// "34 第34回 前半 立花ベース in 初台"

	for _, c := range cases {
		result := ParseTitle(c.Input)

		if c.ExpectedResult.EpTitle != result.EpTitle {
			t.Fatalf("eptitle mismatch: must %s, got %s", c.ExpectedResult.EpTitle, result.EpTitle)
		}

		if c.ExpectedResult.Title != result.Title {
			t.Fatalf("title mismatch: must %s, got %s", c.ExpectedResult.Title, result.Title)
		}

		if c.ExpectedResult.PartIdx != result.PartIdx {
			t.Fatalf("partidx mismatch: must %v, got %v", c.ExpectedResult.PartIdx, result.PartIdx)
		}

	}
}
