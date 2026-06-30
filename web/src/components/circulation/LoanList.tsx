import { useState, useEffect } from "react";
import { getActiveLoans, type Loan } from "../../lib/api";
import LoanTable from "./LoanTable";

export default function LoanList() {
  const [loans, setLoans] = useState<Loan[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setLoading(true);
    setError(null);
    getActiveLoans()
      .then(setLoans)
      .catch((err) => setError(err instanceof Error ? err.message : "Failed to load active loans"))
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return <p class="text-reading-lamp text-body">Loading active loans...</p>;
  }
  if (error) {
    return <p class="text-red-700 text-body bg-red-50 rounded p-3">{error}</p>;
  }

  return <LoanTable loans={loans} />;
}
