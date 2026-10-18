// Crockford Base32 alphabet
const CROCKFORD_ALPHABET = "0123456789abcdefghjkmnpqrstvwxyz";
const SAFE_BASE16_ALPHABET = "0123456789bcdfgh";
const HEX_ALPHABET = "0123456789abcdef";

let counter = 0;

function getMicroTimeWithCounter() {
  const nowMicro = BigInt(Date.now()) * 1000n + BigInt(Math.floor(Math.random() * 1000));
  counter = (counter + 1) & 0xF;
  return (nowMicro << 4n) | BigInt(counter);
}

function encodeBase(val, alphabet) {
  const base = BigInt(alphabet.length);
  if (val === 0n) return alphabet[0];
  let res = "";
  while (val > 0n) {
    const rem = Number(val % base);
    res = alphabet[rem] + res;
    val = val / base;
  }
  return res;
}

function generateShortID() {
  const val = getMicroTimeWithCounter();
  return encodeBase(val, CROCKFORD_ALPHABET).padStart(11, '0');
}

function generateSafeShortID() {
  const val = getMicroTimeWithCounter();
  return encodeBase(val, SAFE_BASE16_ALPHABET).padStart(14, '0');
}

function generateHexShortID() {
  const val = getMicroTimeWithCounter();
  return encodeBase(val, HEX_ALPHABET).padStart(14, '0');
}

function padZero(num, size = 2) {
  return String(num).padStart(size, '0');
}

function getRandomDigits(len) {
  let res = "";
  for (let i = 0; i < len; i++) {
    res += Math.floor(Math.random() * 10);
  }
  return res;
}

function generateHumanUid(formatted) {
  const now = new Date();
  const datePart = now.getUTCFullYear().toString() +
    padZero(now.getUTCMonth() + 1) +
    padZero(now.getUTCDate()) +
    padZero(now.getUTCHours()) +
    padZero(now.getUTCMinutes()) +
    padZero(now.getUTCSeconds()) +
    padZero(now.getUTCMilliseconds(), 3) + "0000"; // 21 digits
  const randPart = getRandomDigits(11); // total 32 digits
  const raw = (datePart + randPart).substring(0, 32);

  if (formatted) {
    return raw.slice(0, 8) + "-" + raw.slice(8, 12) + "-" + raw.slice(12, 16) + "-" + raw.slice(16, 32);
  }
  return raw;
}

function generateMinUid(formatted) {
  const now = new Date();
  const raw = now.getUTCFullYear().toString() +
    padZero(now.getUTCMonth() + 1) +
    padZero(now.getUTCDate()) +
    padZero(now.getUTCHours()) +
    padZero(now.getUTCMinutes());
  if (formatted) {
    return raw.slice(0, 8) + "-" + raw.slice(8, 12);
  }
  return raw;
}

function generateHourUid(formatted) {
  const now = new Date();
  const raw = now.getUTCFullYear().toString() +
    padZero(now.getUTCMonth() + 1) +
    padZero(now.getUTCDate()) +
    padZero(now.getUTCHours());
  if (formatted) {
    return raw.slice(0, 8) + "-" + raw.slice(8, 10);
  }
  return raw;
}

function generateDayUid(formatted) {
  const now = new Date();
  const raw = now.getUTCFullYear().toString() +
    padZero(now.getUTCMonth() + 1) +
    padZero(now.getUTCDate());
  if (formatted) {
    return raw.slice(0, 4) + "-" + raw.slice(4, 6) + "-" + raw.slice(6, 8);
  }
  return raw;
}

function getISOWeekAndYear(date) {
  const d = new Date(Date.UTC(date.getUTCFullYear(), date.getUTCMonth(), date.getUTCDate()));
  const dayNum = d.getUTCDay() || 7;
  d.setUTCDate(d.getUTCDate() + 4 - dayNum);
  const yearStart = new Date(Date.UTC(d.getUTCFullYear(), 0, 1));
  const weekNo = Math.ceil((((d - yearStart) / 86400000) + 1) / 7);
  return { year: d.getUTCFullYear(), week: weekNo };
}

function generateWeekUid(formatted) {
  const now = new Date();
  const { year, week } = getISOWeekAndYear(now);
  const raw = year.toString() + padZero(week);
  if (formatted) {
    return raw.slice(0, 4) + "-" + raw.slice(4, 6);
  }
  return raw;
}

function generateMonthUid(formatted) {
  const now = new Date();
  const raw = now.getUTCFullYear().toString() +
    padZero(now.getUTCMonth() + 1);
  if (formatted) {
    return raw.slice(0, 4) + "-" + raw.slice(4, 6);
  }
  return raw;
}

function getTimeValSec() {
  const now = new Date();
  const y = BigInt(now.getUTCFullYear() % 100);
  const m = BigInt(now.getUTCMonth() + 1);
  const d = BigInt(now.getUTCDate());
  const hh = BigInt(now.getUTCHours());
  const mm = BigInt(now.getUTCMinutes());
  const ss = BigInt(now.getUTCSeconds());
  return y * 10000000000n + m * 100000000n + d * 1000000n + hh * 10000n + mm * 100n + ss;
}

function getTimeValMin() {
  const now = new Date();
  const y = BigInt(now.getUTCFullYear() % 100);
  const m = BigInt(now.getUTCMonth() + 1);
  const d = BigInt(now.getUTCDate());
  const hh = BigInt(now.getUTCHours());
  const mm = BigInt(now.getUTCMinutes());
  return y * 100000000n + m * 1000000n + d * 10000n + hh * 100n + mm;
}

function getTimeValHour() {
  const now = new Date();
  const y = BigInt(now.getUTCFullYear() % 100);
  const m = BigInt(now.getUTCMonth() + 1);
  const d = BigInt(now.getUTCDate());
  const hh = BigInt(now.getUTCHours());
  return y * 1000000n + m * 10000n + d * 100n + hh;
}

function getTimeValDay() {
  const now = new Date();
  const y = BigInt(now.getUTCFullYear() % 100);
  const m = BigInt(now.getUTCMonth() + 1);
  const d = BigInt(now.getUTCDate());
  return y * 10000n + m * 100n + d;
}

function getTimeValWeek() {
  const now = new Date();
  const { year, week } = getISOWeekAndYear(now);
  const y = BigInt(year % 100);
  const w = BigInt(week);
  return y * 100n + w;
}

function getTimeValMonth() {
  const now = new Date();
  const y = BigInt(now.getUTCFullYear() % 100);
  const m = BigInt(now.getUTCMonth() + 1);
  return y * 100n + m;
}

function getTimeValYear() {
  const now = new Date();
  return BigInt(now.getUTCFullYear() % 100);
}

function encodeFixedJS(val, alphabet, width) {
  return encodeBase(val, alphabet).padStart(width, '0');
}

function generateNanoUid(formatted) {
  const now = new Date();
  const datePart = now.getUTCFullYear().toString() +
    padZero(now.getUTCMonth() + 1) +
    padZero(now.getUTCDate()) +
    padZero(now.getUTCHours()) +
    padZero(now.getUTCMinutes()) +
    padZero(now.getUTCSeconds()) +
    padZero(now.getUTCMilliseconds(), 3) + "000"; // 20 digits
  const randPart = getRandomDigits(3);
  const raw = (datePart + randPart).substring(0, 23);

  if (formatted) {
    return raw.slice(0, 8) + "-" + raw.slice(8, 14) + "-" + raw.slice(14, 20) + "-" + raw.slice(20, 23);
  }
  return raw;
}

function generateMicroUid(formatted) {
  const now = new Date();
  const datePart = now.getUTCFullYear().toString() +
    padZero(now.getUTCMonth() + 1) +
    padZero(now.getUTCDate()) +
    padZero(now.getUTCHours()) +
    padZero(now.getUTCMinutes()) +
    padZero(now.getUTCSeconds()) +
    padZero(now.getUTCMilliseconds(), 3) + "000"; // 20 digits
  const raw = datePart.substring(0, 20);

  if (formatted) {
    return raw.slice(0, 8) + "-" + raw.slice(8, 14) + "-" + raw.slice(14, 20);
  }
  return raw;
}

function generateSecUid(formatted) {
  const now = new Date();
  const raw = now.getUTCFullYear().toString() +
    padZero(now.getUTCMonth() + 1) +
    padZero(now.getUTCDate()) +
    padZero(now.getUTCHours()) +
    padZero(now.getUTCMinutes()) +
    padZero(now.getUTCSeconds()); // 14 digits

  if (formatted) {
    return raw.slice(0, 8) + "-" + raw.slice(8, 14);
  }
  return raw;
}

function generateUuidV4(formatted) {
  const hex = Array.from(crypto.getRandomValues(new Uint8Array(16)))
    .map((b, i) => {
      if (i === 6) b = (b & 0x0f) | 0x40; // v4
      if (i === 8) b = (b & 0x3f) | 0x80; // variant
      return padZero(b.toString(16), 2);
    }).join('');

  if (formatted) {
    return hex.slice(0, 8) + "-" + hex.slice(8, 12) + "-" + hex.slice(12, 16) + "-" + hex.slice(16, 20) + "-" + hex.slice(20);
  }
  return hex;
}

function generateUuidV1(formatted) {
  const hex = generateUuidV4(false).slice(0, 28) + "1000"; // representation mockup
  if (formatted) {
    return hex.slice(0, 8) + "-" + hex.slice(8, 12) + "-1" + hex.slice(13, 16) + "-" + hex.slice(16, 20) + "-" + hex.slice(20);
  }
  return hex;
}

function generateUuidV6(formatted) {
  const hex = generateUuidV4(false);
  const v6Hex = hex.slice(0, 12) + "6" + hex.slice(13);
  if (formatted) {
    return v6Hex.slice(0, 8) + "-" + v6Hex.slice(8, 12) + "-" + v6Hex.slice(12, 16) + "-" + v6Hex.slice(16, 20) + "-" + v6Hex.slice(20);
  }
  return v6Hex;
}

function generateUuidV7(formatted) {
  const ms = Date.now();
  const msHex = padZero(ms.toString(16), 12);
  const randHex = Array.from(crypto.getRandomValues(new Uint8Array(10)))
    .map(b => padZero(b.toString(16), 2)).join('');
  const combined = (msHex + "7" + randHex.slice(1)).slice(0, 32);

  if (formatted) {
    return combined.slice(0, 8) + "-" + combined.slice(8, 12) + "-" + combined.slice(12, 16) + "-" + combined.slice(16, 20) + "-" + combined.slice(20);
  }
  return combined;
}

// Number shortening encodings
const BASE16 = "0123456789abcdef";
const BASE32 = "abcdefghijklmnopqrstuvwxyz234567";
const CROCKFORD = "0123456789ABCDEFGHJKMNPQRSTVWXYZ";
const BASE36 = "0123456789abcdefghijklmnopqrstuvwxyz";
const BASE58 = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz";
const BASE62 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz";
const BASE64 = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_";
const ZBASE32 = "ybndrfg8ejkmcpqxot1uwisza3456789";

function shortenNumber(strVal, alphabet) {
  try {
    let val = BigInt(strVal);
    if (val === 0n) return alphabet[0];
    const base = BigInt(alphabet.length);
    let res = "";
    while (val > 0n) {
      const rem = Number(val % base);
      res = alphabet[rem] + res;
      val = val / base;
    }
    return res;
  } catch (e) {
    return "Invalid input";
  }
}

function copyValue(elementId) {
  const el = document.getElementById(elementId);
  if (!el) return;
  navigator.clipboard.writeText(el.innerText);
  const originalText = el.innerText;
  el.style.color = "#38bdf8";
  setTimeout(() => {
    el.style.color = "var(--success-color)";
  }, 300);
}

function updateAll() {
  const formatted = document.getElementById("format-toggle").checked;

  document.getElementById("out-shortid").innerText = generateShortID();
  document.getElementById("out-safeshortid").innerText = generateSafeShortID();
  document.getElementById("out-hexshortid").innerText = generateHexShortID();

  const h = generateHumanUid(formatted);
  document.getElementById("out-humanuid").innerText = h;
  document.getElementById("tag-humanuid").innerText = h.length + " Chars";

  const n = generateNanoUid(formatted);
  document.getElementById("out-nanouid").innerText = n;
  document.getElementById("tag-nanouid").innerText = n.length + " Chars";

  const m = generateMicroUid(formatted);
  document.getElementById("out-microuid").innerText = m;
  document.getElementById("tag-microuid").innerText = m.length + " Chars";

  const s = generateSecUid(formatted);
  document.getElementById("out-secuid").innerText = s;
  document.getElementById("tag-secuid").innerText = s.length + " Chars";

  const minUid = generateMinUid(formatted);
  if (document.getElementById("out-minuid")) {
    document.getElementById("out-minuid").innerText = minUid;
    document.getElementById("tag-minuid").innerText = minUid.length + " Chars";
  }

  const hourUid = generateHourUid(formatted);
  if (document.getElementById("out-houruid")) {
    document.getElementById("out-houruid").innerText = hourUid;
    document.getElementById("tag-houruid").innerText = hourUid.length + " Chars";
  }

  const dayUid = generateDayUid(formatted);
  if (document.getElementById("out-dayuid")) {
    document.getElementById("out-dayuid").innerText = dayUid;
    document.getElementById("tag-dayuid").innerText = dayUid.length + " Chars";
  }

  const weekUid = generateWeekUid(formatted);
  if (document.getElementById("out-weekuid")) {
    document.getElementById("out-weekuid").innerText = weekUid;
    document.getElementById("tag-weekuid").innerText = weekUid.length + " Chars";
  }

  const monthUid = generateMonthUid(formatted);
  if (document.getElementById("out-monthuid")) {
    document.getElementById("out-monthuid").innerText = monthUid;
    document.getElementById("tag-monthuid").innerText = monthUid.length + " Chars";
  }

  // Precision ShortIDs
  if (document.getElementById("out-shortid-sec")) {
    document.getElementById("out-shortid-sec").innerText = encodeFixedJS(getTimeValSec(), CROCKFORD_ALPHABET, 8);
    document.getElementById("out-safeshortid-sec").innerText = encodeFixedJS(getTimeValSec(), SAFE_BASE16_ALPHABET, 10);
    document.getElementById("out-hexshortid-sec").innerText = encodeFixedJS(getTimeValSec(), HEX_ALPHABET, 10);

    document.getElementById("out-shortid-min").innerText = encodeFixedJS(getTimeValMin(), CROCKFORD_ALPHABET, 7);
    document.getElementById("out-safeshortid-min").innerText = encodeFixedJS(getTimeValMin(), SAFE_BASE16_ALPHABET, 9);
    document.getElementById("out-hexshortid-min").innerText = encodeFixedJS(getTimeValMin(), HEX_ALPHABET, 9);

    document.getElementById("out-shortid-hour").innerText = encodeFixedJS(getTimeValHour(), CROCKFORD_ALPHABET, 6);
    document.getElementById("out-safeshortid-hour").innerText = encodeFixedJS(getTimeValHour(), SAFE_BASE16_ALPHABET, 7);
    document.getElementById("out-hexshortid-hour").innerText = encodeFixedJS(getTimeValHour(), HEX_ALPHABET, 7);

    document.getElementById("out-shortid-day").innerText = encodeFixedJS(getTimeValDay(), CROCKFORD_ALPHABET, 4);
    document.getElementById("out-safeshortid-day").innerText = encodeFixedJS(getTimeValDay(), SAFE_BASE16_ALPHABET, 5);
    document.getElementById("out-hexshortid-day").innerText = encodeFixedJS(getTimeValDay(), HEX_ALPHABET, 5);

    document.getElementById("out-shortid-week").innerText = encodeFixedJS(getTimeValWeek(), CROCKFORD_ALPHABET, 3);
    document.getElementById("out-safeshortid-week").innerText = encodeFixedJS(getTimeValWeek(), SAFE_BASE16_ALPHABET, 4);
    document.getElementById("out-hexshortid-week").innerText = encodeFixedJS(getTimeValWeek(), HEX_ALPHABET, 4);

    document.getElementById("out-shortid-month").innerText = encodeFixedJS(getTimeValMonth(), CROCKFORD_ALPHABET, 3);
    document.getElementById("out-safeshortid-month").innerText = encodeFixedJS(getTimeValMonth(), SAFE_BASE16_ALPHABET, 4);
    document.getElementById("out-hexshortid-month").innerText = encodeFixedJS(getTimeValMonth(), HEX_ALPHABET, 4);

    document.getElementById("out-shortid-year").innerText = encodeFixedJS(getTimeValYear(), CROCKFORD_ALPHABET, 2);
    document.getElementById("out-safeshortid-year").innerText = encodeFixedJS(getTimeValYear(), SAFE_BASE16_ALPHABET, 2);
    document.getElementById("out-hexshortid-year").innerText = encodeFixedJS(getTimeValYear(), HEX_ALPHABET, 2);
  }

  const u4 = generateUuidV4(formatted);
  document.getElementById("out-uuidv4").innerText = u4;
  document.getElementById("tag-uuidv4").innerText = u4.length + " Chars";

  const u1 = generateUuidV1(formatted);
  document.getElementById("out-uuidv1").innerText = u1;
  document.getElementById("tag-uuidv1").innerText = u1.length + " Chars";

  const u6 = generateUuidV6(formatted);
  document.getElementById("out-uuidv6").innerText = u6;
  document.getElementById("tag-uuidv6").innerText = u6.length + " Chars";

  const u7 = generateUuidV7(formatted);
  document.getElementById("out-uuidv7").innerText = u7;
  document.getElementById("tag-uuidv7").innerText = u7.length + " Chars";

  const nowSec = Math.floor(Date.now() / 1000).toString();
  const nowMicro = (BigInt(Date.now()) * 1000n).toString();
  const nowNano = (BigInt(Date.now()) * 1000000n).toString();

  document.getElementById("out-timestamp").innerText = nowSec;
  document.getElementById("out-timestamp-micro").innerText = nowMicro;
  document.getElementById("out-timestamp-nano").innerText = nowNano;

  updateShortenerResults();
}

function updateShortenerResults() {
  const inputVal = document.getElementById("shorten-input").value.trim();
  const resultsContainer = document.getElementById("shortener-results");

  const encodings = [
    { name: "Base16", alphabet: BASE16 },
    { name: "Base32", alphabet: BASE32 },
    { name: "Crockford", alphabet: CROCKFORD },
    { name: "z-base-32", alphabet: ZBASE32 },
    { name: "Base36", alphabet: BASE36 },
    { name: "Base58", alphabet: BASE58 },
    { name: "Base62", alphabet: BASE62 },
    { name: "Base64", alphabet: BASE64 }
  ];

  resultsContainer.innerHTML = encodings.map(enc => {
    const shortVal = shortenNumber(inputVal, enc.alphabet);
    return `
      <div class="card" style="padding: 0.9rem;">
        <div class="card-header" style="margin-bottom: 0.4rem;">
          <span class="card-title" style="font-size: 0.95rem;">${enc.name}</span>
          <span class="meta-tag">${shortVal.length} Chars</span>
        </div>
        <div class="output-box" style="padding: 0.5rem;">
          <span class="output-value" style="font-size: 0.85rem;">${shortVal}</span>
        </div>
      </div>
    `;
  }).join('');
}

document.addEventListener("DOMContentLoaded", () => {
  document.getElementById("generate-all-btn").addEventListener("click", updateAll);
  document.getElementById("format-toggle").addEventListener("change", updateAll);
  document.getElementById("shorten-input").addEventListener("input", updateShortenerResults);
  document.getElementById("random-shorten-input-btn").addEventListener("click", () => {
    document.getElementById("shorten-input").value = generateMicroUid(false);
    updateShortenerResults();
  });

  // Initial render
  updateAll();
});
