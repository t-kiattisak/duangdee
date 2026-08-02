# Database Design Analysis: PostgreSQL TEXT[] vs Normalized Relational Table

This document compares the architectural trade-offs between storing card keywords using **PostgreSQL Native Array (`TEXT[]`)** vs a **Normalized Relational Table (`keywords` & `card_meaning_keywords`)**.

---

## 1. Comparative Analysis Matrix

| Comparison Dimension | Option A: Native Array (`keywords TEXT[]`) | Option B: Normalized Junction Table |
| :--- | :--- | :--- |
| **Database Structure** | Single column inside `card_meanings` table. | Requires 3 tables: `card_meanings`, `keywords`, `card_meaning_keywords`. |
| **Query Complexity** | 🟢 **Extremely Simple**: `SELECT * FROM card_meanings` (Zero JOINs). | 🔴 **Complex**: Requires 2 `INNER JOIN` operations per query. |
| **Read Performance (Latency)** | ⚡ **Ultra Fast (< 1ms)**: Fetches all keywords in a single row read. | 🐢 **Slower**: Multi-row index lookups and join processing. |
| **Redis Caching Compatibility**| 🟢 **Direct JSON Serialization**: Maps 1-to-1 with Go `[]string` and Redis JSON. | 🟡 **Requires Assembly**: Must assemble JOIN results into array before caching. |
| **Dynamic Tag Management** | 🔴 Needs string matching within array (`ANY(keywords)`). | 🟢 **Easier**: Rename a global keyword in one single row update. |

---

## 2. Why Option A (`TEXT[]`) is Superior for Tarot Knowledge Base

1. **Tarot Meanings are Read-Heavy & Static**:
   - Tarot card meanings and keywords rarely change once inserted.
   - We read card meanings thousands of times per second, but write/update them almost zero times after initial setup.
2. **Zero-JOIN Query Speed**:
   - Reading Service needs to query meanings for 3 to 10 cards simultaneously in a reading session.
   - With `TEXT[]`, a single query retrieves everything in 1 round-trip without JOIN overhead.
3. **Native PostgreSQL Support**:
   - PostgreSQL is unique among SQL databases in offering first-class, indexed array types (`ANY(keywords)`, GIN indexing for array search).

---

## 3. When Option B (Normalized Junction Table) Would Be Preferred

A normalized junction table would be chosen IF:
- Keywords were dynamic user-generated tags (e.g. Medium/Twitter hashtags where users create millions of tags).
- We needed strict Foreign Key constraints on keyword names to enforce spelling uniqueness across millions of records.
- Keywords had their own metadata (e.g. `keyword_id`, `created_by`, `popularity_score`).

Since Tarot keywords are **static dictionary metadata owned by our system**, **Option A (`TEXT[]`)** delivers maximum performance and simplicity.
