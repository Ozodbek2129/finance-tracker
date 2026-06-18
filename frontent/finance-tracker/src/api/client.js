// ── BASE CONFIGURATION ──────────────────────────────────────────────────────
const BASE_URL = "https://finance-tracker-production-1ff6.up.railway.app/api/v1";

async function request(method, path, body = null, params = {}) {
  const url = new URL(`${BASE_URL}${path}`);
  Object.entries(params).forEach(([k, v]) => {
    if (v !== undefined && v !== null) url.searchParams.set(k, v);
  });

  const options = {
    method,
    headers: { "Content-Type": "application/json" },
  };
  if (body) options.body = JSON.stringify(body);

  const res = await fetch(url.toString(), options);
  if (!res.ok) {
    const err = await res.json().catch(() => ({}));
    throw new Error(err.message || `Xato: ${res.status}`);
  }
  return res.json();
}

// ── API OBJECTS ─────────────────────────────────────────────────────────────
export const categoriesAPI = {
  getAll: (params = {}) => request("GET", "/categories", null, params).then(r => r.categories || []),
  create: (data) => request("POST", "/categories", data),
  update: (id, data) => request("PUT", `/categories/${id}`, data),
  delete: (id) => request("DELETE", `/categories/${id}`),
};

export const transactionsAPI = {
  getAll: (params = {}) => request("GET", "/transactions", null, params).then(r => r.transactions || []),
  create: (data) => request("POST", "/transactions", data),
  update: (id, data) => request("PUT", `/transactions/${id}`, data),
  delete: (id) => request("DELETE", `/transactions/${id}`),
};

export const statsAPI = {
  get: (params = {}) => request("GET", "/stats", null, params),
};

// ── ANIMATION ENGINE (COUNT UP) ─────────────────────────────────────────────
/**
 * Raqamlarni noldan boshlab maqsadli songacha silliq o'stiruvchi animatsiya
 * @param {string} elementId - HTML elementning ID si
 * @param {number} targetValue - Yakuniy hisob kitob summasi
 */
function animateCount(elementId, targetValue) {
  const element = document.getElementById(elementId);
  if (!element) return; // Agar element HTML ichida topilmasa, xato bermaydi

  const duration = 1200; // Animatsiya davomiyligi: 1.2 sekund
  const startTime = performance.now();

  function updateNumber(currentTime) {
    const elapsedTime = currentTime - startTime;

    if (elapsedTime >= duration) {
      element.innerText = formatMoney(targetValue);
      return;
    }

    const progress = elapsedTime / duration;
    // Ease-out effekti: boshida tez o'sadi, oxirida sekinlashadi
    const easeOutProgress = 1 - Math.pow(1 - progress, 3);
    const currentValue = Math.floor(easeOutProgress * targetValue);

    element.innerText = formatMoney(currentValue);
    requestAnimationFrame(updateNumber);
  }

  requestAnimationFrame(updateNumber);
}

/**
 * Raqamni pul formatiga o'tkazish (masalan: 4000000 -> "4 000 000")
 */
function formatMoney(amount) {
  return amount.toString().replace(/\B(?=(\d{3})+(?!\d))/g, " ");
}

// ── DATA CORE (BACKEND INTEGRATION) ─────────────────────────────────────────
/**
 * Tanlangan oy va yil bo'yicha statistikani backend'dan oladi va animatsiya qiladi
 * @param {string} month - Masalan: "Iyun" yoki raqamda "06" (backend talabiga qarab)
 * @param {number|string} year - Masalan: 2026
 */
export async function loadAndAnimateStatistics(month, year) {
  try {
    // 1. Backend API ga so'rov yuboramiz
    const stats = await statsAPI.get({ month, year });

    // 2. Kelgan ma'lumotlarni o'zgaruvchilarga olamiz 
    // (Agar kelayotgan kalit so'zlar boshqacha bo'lsa, shu yerda nomini to'g'rilaysiz)
    const daromad = stats.total_income || stats.income || 0;
    const xarajat = stats.total_expense || stats.expense || 0;
    const balans = stats.balance !== undefined ? stats.balance : (daromad - xarajat);

    // 3. HTML elementlardagi raqamlarni animatsiya bilan yangilaymiz
    animateCount('daromad-val', daromad);
    animateCount('xarajat-val', xarajat);
    animateCount('balans-val', balans);

  } catch (error) {
    console.error("Statistikani yuklash va animatsiya qilishda xatolik:", error);
  }
}

// ── INITIALIZATION (SAHIFA YUKLANIShI) ──────────────────────────────────────
// Sahifa birinchi marta ochilganda avtomatik ishga tushirish qismi
window.addEventListener('DOMContentLoaded', () => {
  // Boshlanishiga joriy oy va yil ma'lumotlarini yuklaymiz (Iyun, 2026)
  loadAndAnimateStatistics("Iyun", 2026);
});