package queries

// INSERT NEW URL
const InsertURL string = `INSERT INTO urls (url, short_url) VALUES (?,?)`

// Find URL by Short URL
const FindByShort string = `SELECT id, url, short_url, click_count, created_at FROM urls WHERE short_url = ?`

// Find All Url
const FindAll string = `SELECT id, url, short_url, click_count, created_at FROM urls`

// Delete URL by ID
const DeleteByID = `DELETE FROM urls WHERE id = ?`
