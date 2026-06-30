const BASE_URL = "/api/v1";
const TOKEN_KEY = "lmss_token";

interface ApiResponse<T> {
  success: boolean;
  data: T;
  message: string;
}

function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY);
}

export function setToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token);
}

export function clearToken(): void {
  localStorage.removeItem(TOKEN_KEY);
}

export function isAuthenticated(): boolean {
  return !!getToken();
}

async function request<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const token = getToken();
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(options.headers as Record<string, string>),
  };

  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const response = await fetch(`${BASE_URL}${endpoint}`, {
    ...options,
    headers,
  });

  if (response.status === 401 && endpoint !== "/auth/login") {
    clearToken();
    window.location.href = "/login";
    throw new Error("Session expired");
  }

  const json: ApiResponse<T> = await response.json();

  if (!response.ok || !json.success) {
    throw new Error(json.message || "Request failed");
  }

  return json.data;
}

export const api = {
  get: <T>(endpoint: string) => request<T>(endpoint),
  post: <T>(endpoint: string, data: unknown) =>
    request<T>(endpoint, {
      method: "POST",
      body: JSON.stringify(data),
    }),
  put: <T>(endpoint: string, data: unknown) =>
    request<T>(endpoint, {
      method: "PUT",
      body: JSON.stringify(data),
    }),
  delete: <T>(endpoint: string) =>
    request<T>(endpoint, {
      method: "DELETE",
    }),
};

export interface LoginResponse {
  token: string;
  user: {
    id: string;
    email: string;
    name: string;
    role: string;
  };
}

export interface OverviewStats {
  books_checked_out_today: number;
  new_patrons_this_week: number;
  overdue_returns: number;
  recent_activity: ActivityEntry[];
}

export interface ActivityEntry {
  book_title: string;
  patron_name: string;
  action: string;
  timestamp: string;
}

export async function login(email: string, password: string): Promise<LoginResponse> {
  const data = await api.post<LoginResponse>("/auth/login", { email, password });
  setToken(data.token);
  return data;
}

export async function logout(): Promise<void> {
  try {
    await api.post("/auth/logout", {});
  } finally {
    clearToken();
  }
}

export interface Book {
  id: string;
  title: string;
  author: string;
  isbn?: string;
  genre?: string;
  status: string;
  location?: string;
  created_at: string;
  updated_at: string;
}

export async function listBooks(): Promise<Book[]> {
  return api.get<Book[]>("/catalog");
}

export async function searchBooks(query: string): Promise<Book[]> {
  return api.get<Book[]>(`/catalog/search?q=${encodeURIComponent(query)}`);
}

export async function getBook(id: string): Promise<Book> {
  return api.get<Book>(`/catalog/${id}`);
}

export async function createBook(data: {
  title: string;
  author: string;
  isbn?: string;
  genre?: string;
  location?: string;
  status: string;
}): Promise<Book> {
  return api.post<Book>("/catalog", data);
}

export async function deleteBook(id: string): Promise<void> {
  await api.delete(`/catalog/${id}`);
}

export interface Patron {
  id: string;
  name: string;
  email: string;
  phone?: string;
  active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Loan {
  id: string;
  book_id: string;
  book_title: string;
  patron_id: string;
  patron_name: string;
  checked_out_at: string;
  due_date: string;
  returned_at?: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export async function listPatrons(): Promise<Patron[]> {
  return api.get<Patron[]>("/patrons");
}

export async function searchPatrons(query: string): Promise<Patron[]> {
  return api.get<Patron[]>(`/patrons/search?q=${encodeURIComponent(query)}`);
}

export async function getPatron(id: string): Promise<Patron> {
  return api.get<Patron>(`/patrons/${id}`);
}

export async function createPatron(data: {
  name: string;
  email: string;
  phone?: string;
  active?: boolean;
}): Promise<Patron> {
  return api.post<Patron>("/patrons", data);
}

export async function updatePatron(id: string, data: {
  name?: string;
  email?: string;
  phone?: string;
  active?: boolean;
}): Promise<void> {
  await api.put(`/patrons/${id}`, data);
}

export async function deletePatron(id: string): Promise<void> {
  await api.delete(`/patrons/${id}`);
}

export async function checkoutBook(data: {
  book_id: string;
  patron_id: string;
}): Promise<Loan> {
  return api.post<Loan>("/circulation/checkout", data);
}

export async function returnBook(loanId: string): Promise<void> {
  await api.post("/circulation/returns", { loan_id: loanId });
}

export async function getActiveLoans(): Promise<Loan[]> {
  return api.get<Loan[]>("/circulation/loans/active");
}

export async function getLoanHistory(): Promise<Loan[]> {
  return api.get<Loan[]>("/circulation/loans/history");
}

export interface Fine {
  id: string;
  patron_id: string;
  patron_name: string;
  book_title: string;
  amount: number;
  status: string;
  paid: boolean;
  due_date: string;
  paid_at?: string;
  created_at: string;
  updated_at: string;
}

export interface Reservation {
  id: string;
  book_id: string;
  book_title: string;
  patron_id: string;
  patron_name: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export async function listFines(): Promise<Fine[]> {
  return api.get<Fine[]>("/fines");
}

export async function payFine(fineId: string): Promise<Fine> {
  return api.post<Fine>("/fines/pay", { fine_id: fineId });
}

export async function listReservations(): Promise<Reservation[]> {
  return api.get<Reservation[]>("/reservations");
}

export async function createReservation(data: {
  book_id: string;
  patron_id: string;
}): Promise<Reservation> {
  return api.post<Reservation>("/reservations", data);
}

export async function deleteReservation(id: string): Promise<void> {
  await api.delete(`/reservations/${id}`);
}

export async function getDashboardOverview(): Promise<OverviewStats> {
  return api.get<OverviewStats>("/reports/overview");
}
