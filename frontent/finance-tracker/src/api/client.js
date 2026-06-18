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

// ── Categories ──────────────────────────────────────────────────────────────
export const categoriesAPI = {
  getAll: (params = {}) => request("GET", "/categories", null, params).then(r => r.categories || []),
  create: (data) => request("POST", "/categories", data),
  update: (id, data) => request("PUT", `/categories/${id}`, data),
  delete: (id) => request("DELETE", `/categories/${id}`),
};

// ── Transactions ─────────────────────────────────────────────────────────────
export const transactionsAPI = {
  getAll: (params = {}) => request("GET", "/transactions", null, params).then(r => r.transactions || []),
  create: (data) => request("POST", "/transactions", data),
  update: (id, data) => request("PUT", `/transactions/${id}`, data),
  delete: (id) => request("DELETE", `/transactions/${id}`),
};

// ── Stats ─────────────────────────────────────────────────────────────────────
export const statsAPI = {
  get: (params = {}) => request("GET", "/stats", null, params),
};