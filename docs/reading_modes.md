# Tarot Reading Modes Specification

## 1. Overview of Reading Modes

To accommodate diverse user needs, the system supports **5 Core Tarot Reading Modes**, categorized by **Time Horizons** and **Life Contexts**. Each mode utilizes a specific Spread Layout, Time Horizon, and Card Meaning Dictionary Context.

---

## 2. Reading Modes Breakdown Table

| Mode Name (ชื่อโหมด) | Target Time Horizon | Spread Layout Used | Number of Cards | Coin Cost | User Intent / Purpose |
| :--- | :--- | :--- | :--- | :--- | :--- |
| 1. **Daily Horoscope (ดวงรายวัน)** | 1 Day (วันนี้) | `single-card` | 1 Card | **0 Coins (Free)** / 5 Coins | แนวทาง/คำเตือนสั้นๆ ประจำวัน |
| 2. **Monthly Forecast (ดวงรายเดือน)** | 1 Month (เดือนนี้) | `three-card` or `four-card` | 3 - 4 Cards | **10 Coins** | ภาพรวมการงาน การเงิน ความรัก ในเดือนนี้ |
| 3. **Yearly Forecast (ดวงรายปี)** | 1 Year (ปีนี้/ปีหน้า) | `twelve-card-zodiac` | 12 Cards | **50 Coins** | ภาพรวมทั้ง 12 เดือนประจำปีตามจักรราศี |
| 4. **Specific Question (ถามเรื่องเฉพาะเจาะจง)**| Custom / Immediate | `three-card` / `celtic-cross` | 3 - 10 Cards | **10 - 30 Coins** | ถามเฉพาะเรื่อง (งาน, เงิน, ความรัก, สุขภาพ) |
| 5. **Yes/No Quick Answer (ตอบคำถาม ใช่/ไม่ใช่)**| Immediate | `single-card-yesno` | 1 Card | **5 Coins** | ต้องการคำตอบรวดเร็วว่าจะทำสิ่งนั้นดีหรือไม่ |

---

## 3. Detailed Breakdown of Each Reading Mode

### Mode 1: Daily Horoscope (ดวงประจำวัน - วันนี้)
- **User Intent**: ผู้ใช้ต้องการทราบแนวทางการดำเนินชีวิต ข้อควรระวัง หรือสติประจำวัน
- **Card Mechanics**: สุ่มเปิดไพ่ 1 ใบ (หัวตั้ง/หัวกลับ)
- **Interpretation Logic**: ดึงความหมายจากหมวด `general` + `upright/reversed` สรุปคำทำนายสั้น 2-3 บรรทัด

### Mode 2: Monthly Forecast (ดวงประจำเดือน - เดือนนี้)
- **User Intent**: เช็กแนวโน้มชีวิตตลอดทั้งเดือน แบ่งตามหมวดหมู่หลัก
- **Card Mechanics**: เปิดไพ่ 4 ใบ
  - ใบที่ 1: ภาพรวมประจำเดือน (Overall Theme)
  - ใบที่ 2: ดวงการงาน & อุปสรรค (Career)
  - ใบที่ 3: ดวงการเงิน &โชคลาภ (Finance)
  - ใบที่ 4: ดวงความรัก & ความสัมพันธ์ (Love)
- **Interpretation Logic**: รวบรวมความหมายรายใบตามหมวดหมู่ และสรุปภาพรวมประจำเดือน

### Mode 3: Yearly Forecast (ดวงประจำปี - 12 เดือน)
- **User Intent**: มองภาพใหญ่ชีวิตตลอดทั้งปี เหมาะสำหรับการดูช่วงปีใหม่ หรือวันเกิด
- **Card Mechanics**: เปิดไพ่ 12 ใบ จัดวางวงกลมตามจักรราศี (Zodiac Wheel) 12 เรือน (House 1 ถึง House 12)
  - เรือนที่ 1 (มกราคม/ตัวตน) -> เรือนที่ 12 (ธันวาคม/สิ่งเร้นลับ/บทสรุปปี)
- **Interpretation Logic**: แสดงผลกราฟคำทำนายรายเดือนตลอดทั้งปี

### Mode 4: Specific Question Focus (ถามเรื่องเฉพาะเจาะจง)
- **User Intent**: มีเรื่องกลุ้มใจหรือคำถามในใจชัดเจน (เช่น "เขาคิดยังไงกับเรา", "จะได้เปลี่ยนงานไหม")
- **Card Mechanics**:
  - ผู้ใช้เลือก **หมวดหมู่** (ความรัก, การงาน, การเงิน, สุขภาพ)
  - ผู้ใช้พิมพ์ **คำถาม** หรือเลือกคำถามยอดนิยม
  - เลือก Spread (3 ใบ: อดีต-ปัจจุบัน-อนาคต หรือ 10 ใบ Celtic Cross)
- **Interpretation Logic**: แมตช์ความหมายไพ่ตามหมวดหมู่ที่ผู้ใช้เลือก (e.g. ดึงเฉพาะความหมายหมวด `love`)

### Mode 5: Quick Yes / No (ตอบคำถามสั้น ใช่/ไม่ใช่)
- **User Intent**: ตัดสินใจเรื่องด่วน ต้องการคำตอบกระชับ
- **Card Mechanics**: สุ่มเปิด 1 ใบ
- **Interpretation Logic**: ให้คะแนนประเภทไพ่ (Positive/Negative/Neutral) + Orientation (หัวตั้ง = Yes, หัวกลับ = No)

---

## 4. User Pre-Reading Intention / Question Prompt (ข้อ 1: สิ่งที่ User ต้องการก่อนจับไพ่)

ตอบคำถาม: **"1 เป็นแนวแบบว่าสิ่งที่ user ต้องการก่อนที่จะจับไพ่ใช่ป่ะ"**
-> **ถูกต้องเลยครับ!**

ก่อนที่ผู้ใช้จะเริ่มกดสับหรือจับไพ่ ระบบจะให้ผู้ใช้ทำ **"Pre-Reading Intention Setting" (ตั้งจิต/ระบุความตั้งใจ)** เพื่อให้การทำนายแม่นยำตรงจุด:

```
[ Step 1: เลือกโหมด ] ---> [ Step 2: เลือกหมวด / พิมพ์คำถาม ] ---> [ Step 3: ตั้งจิตและกดจับไพ่ ]
(ดวงวันนี้ / เดือนนี้ / ปีนี้)  (งาน / เงิน / ความรัก / พิมพ์ถาม)     (สุ่มไพ่ด้วย Crypto Randomizer)
```

1. **เลือกโหมดการดู**: วันนี้ / เดือนนี้ / ปีนี้ / ถามเฉพาะเรื่อง
2. **ระบุเรื่องที่กังวล (Topic Context)**: การงาน / การเงิน / ความรัก / สุขภาพ
3. **พิมพ์คำถามตั้งจิต (Optional Question Prompt)**: เช่น *"ควรลาออกจากงานตอนนี้ไหม?"* (คำถามนี้จะถูกบันทึกไว้ใน `reading_sessions.question` เพื่อให้นำไปแสดงคู่กับผลไพ่ และใช้เก็บใน Journal)
