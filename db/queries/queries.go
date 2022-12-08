package queries

// INSERT NEW URL
const InsertURL string = `INSERT INTO urls (url, short_url, created_at, user_id) VALUES (?,?,?,?)`

// Find URL by Short URL
const FindByShort string = `SELECT id, url, short_url, click_count, created_at FROM urls WHERE short_url = ?`

// Find URL by URL ID
const FindByID string = `SELECT id, url, short_url, click_count, created_at FROM urls WHERE id = ?`

// Find All Url
const FindAll string = `SELECT id, url, short_url, click_count, created_at FROM urls`

// Delete URL by ID
const DeleteByID = `DELETE FROM urls WHERE id = ?`

//Increment click_count when url opened
const IncrementClickCount = `UPDATE urls SET click_count = ? WHERE id = ?`
