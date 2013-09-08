/*****************************************************************************
 **
 ** Gofr
 ** https://github.com/melllvar/Gofr
 ** Copyright (C) 2013 Akop Karapetyan
 **
 ** This program is free software; you can redistribute it and/or modify
 ** it under the terms of the GNU General Public License as published by
 ** the Free Software Foundation; either version 2 of the License, or
 ** (at your option) any later version.
 **
 ** This program is distributed in the hope that it will be useful,
 ** but WITHOUT ANY WARRANTY; without even the implied warranty of
 ** MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 ** GNU General Public License for more details.
 **
 ** You should have received a copy of the GNU General Public License
 ** along with this program; if not, write to the Free Software
 ** Foundation, Inc., 675 Mass Ave, Cambridge, MA 02139, USA.
 **
 ******************************************************************************
 */
 
package storage

import (
  "appengine/datastore"
  "html"
  "rss"
  "regexp"
  "sanitize"
  "time"
)

type Feed struct {
  URL string
  Title string
  Description string `datastore:",noindex"`
  Updated time.Time
  Link string
  Format string
  Retrieved time.Time
  HourlyUpdateFrequency float32
  UpdateCounter int64
}

type EntryMeta struct {
  Published time.Time
  Retrieved time.Time
  Updated time.Time
  UpdateOffset int64
  Entry *datastore.Key
}

type Entry struct {
  UniqueID string     `json:"-"`
  Author string       `json:"author"`
  Title string        `json:"title"`
  Link string         `json:"link"`
  Published time.Time `json:"published"`
  Updated time.Time   `json:"-"`

  Content string `json:"content" datastore:",noindex"`
  Summary string `json:"summary"`
}

type UserSubscriptions struct {
  Subscriptions []Subscription `json:"subscriptions"`
  Folders []Folder             `json:"folders"`
}

type FolderRef struct {
  UserID UserID
  FolderID string
}

type SubscriptionRef struct {
  FolderRef
  SubscriptionID string
}

type Subscription struct {
  ID string     `datastore:"-" json:"id"`
  Link string   `datastore:"-" json:"link"`
  Parent string `datastore:"-" json:"parent,omitempty"`

  Updated time.Time    `json:"-"`
  Subscribed time.Time `json:"-"`
  Feed *datastore.Key  `json:"-"`

  Title string         `json:"title"`
  UnreadCount int      `json:"unread"`
}

type ArticlePage struct {
  Articles []Article `json:"articles"`
  Continue string    `json:"continue,omitempty"`
}

type ArticleRef struct {
  SubscriptionRef
  ArticleID string
}

type ArticleFilter struct {
  UserID UserID
  FolderID string
  SubscriptionID string
  Property string
}

type Article struct {
  ID string             `datastore:"-" json:"id"`
  Source string         `datastore:"-" json:"source"`
  Details *Entry        `datastore:"-" json:"details"`

  Published time.Time   `json:"-"`
  Retrieved time.Time   `json:"-"`
  Entry *datastore.Key  `json:"-"`

  Properties []string   `json:"properties"`
}

type Folder struct {
  ID string    `datastore:"-" json:"id"`
  Title string `json:"title"`
}

var extraSpaceStripper *regexp.Regexp = regexp.MustCompile(`\s\s+`)

func generateSummary(entry *rss.Entry) string {
  // TODO: This process should be streamlined to do
  // more things with fewer passes
  sanitized := sanitize.StripTags(entry.Content)
  unescaped := html.UnescapeString(sanitized)
  stripped := extraSpaceStripper.ReplaceAllString(unescaped, "")

  if runes := []rune(stripped); len(runes) > 400 {
    return string(runes[:400])
  } else {
    return stripped
  }
}