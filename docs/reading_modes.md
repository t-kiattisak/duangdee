# Tarot Reading Modes Specification

## 1. Overview of Reading Modes

To accommodate diverse user needs, the system supports **5 Core Tarot Reading Modes**, categorized by **Time Horizons** and **Life Contexts**. Each mode utilizes a specific Spread Layout, Time Horizon, and Card Meaning Dictionary Context.

---

## 2. Reading Modes Breakdown Table

| Mode Name | Target Time Horizon | Spread Layout Used | Number of Cards | Coin Cost | User Intent / Purpose |
| :--- | :--- | :--- | :--- | :--- | :--- |
| 1. **Daily Horoscope** | 1 Day (Today) | `single-card` | 1 Card | **0 Coins (Daily Free)** / 5 Coins | Daily guidance, mindfulness & key advice. |
| 2. **Monthly Forecast** | 1 Month (This Month) | `three-card` or `four-card` | 3 - 4 Cards | **10 Coins** | Overall theme + Career, Finance, Love forecast. |
| 3. **Yearly Forecast** | 1 Year (This Year) | `twelve-card-zodiac` | 12 Cards | **50 Coins** | Full 12-month zodiac wheel forecast. |
| 4. **Specific Question Focus**| Custom / Immediate | `three-card` / `celtic-cross` | 3 - 10 Cards | **10 - 30 Coins** | Dedicated inquiry (Love, Career, Finance, Health). |
| 5. **Quick Yes / No Answer**| Immediate | `single-card-yesno` | 1 Card | **5 Coins** | Quick binary decision & positive/negative indicator. |

---

## 3. Detailed Breakdown of Each Reading Mode

### Mode 1: Daily Horoscope (Today)
- **User Intent**: Understand daily mindfulness, warnings, or key guidance for today.
- **Card Mechanics**: Draw 1 card (Upright / Reversed).
- **Interpretation Logic**: Fetches card dictionary under `general` context + `upright/reversed` orientation. Generates a concise 2-3 line daily summary.

### Mode 2: Monthly Forecast (This Month)
- **User Intent**: Evaluate life trends across the month categorized by key life aspects.
- **Card Mechanics**: Draw 4 cards:
  - Card 1: Monthly Overall Theme.
  - Card 2: Career & Obstacles.
  - Card 3: Finance & Luck.
  - Card 4: Love & Relationships.
- **Interpretation Logic**: Synthesizes position-by-position meanings and computes a monthly trend summary.

### Mode 3: Yearly Forecast (12 Months)
- **User Intent**: High-level annual perspective, ideal for New Year or Birthday readings.
- **Card Mechanics**: Draw 12 cards arranged in a 12-House Zodiac Wheel:
  - House 1 (January / Self) -> House 12 (December / Closure / Outcome).
- **Interpretation Logic**: Renders a 12-month prediction timeline chart.

### Mode 4: Specific Question Focus (Dedicated Inquiry)
- **User Intent**: Clear, specific query in mind (e.g., "What are their feelings toward me?", "Should I change jobs?").
- **Card Mechanics**:
  - User selects **Category** (Love, Career, Finance, Health).
  - User inputs a **custom question** or selects popular prompt templates.
  - User selects Spread Layout (3 Cards: Past-Present-Future or 10 Cards: Celtic Cross).
- **Interpretation Logic**: Matches dictionary entries specific to the selected category (e.g. fetches `love` category entries).

### Mode 5: Quick Yes / No (Fast Decision)
- **User Intent**: Rapid decision making requiring a direct, binary recommendation.
- **Card Mechanics**: Draw 1 card.
- **Interpretation Logic**: Calculates card archetype polarity (Positive/Negative/Neutral) + Orientation (Upright = Yes, Reversed = No).

---

## 4. User Pre-Reading Intention & Question Prompt Workflow

Before drawing cards, the frontend guides the user through **Pre-Reading Intention Setting**:

```
[ Step 1: Select Mode ]  --->  [ Step 2: Topic & Custom Question ]  --->  [ Step 3: Draw Cards ]
(Daily / Monthly / Yearly)      (Love / Career + Custom Query Prompt)      (CSPRNG Crypto Randomizer)
```

1. **Select Reading Mode**: Daily / Monthly / Yearly / Specific Question.
2. **Specify Topic Context**: Career / Finance / Love / Health.
3. **Input Custom Question Prompt**: e.g., *"Should I accept the new job offer this month?"* (Saved in `reading_sessions.question` to render alongside drawn card results and store in user journals).
