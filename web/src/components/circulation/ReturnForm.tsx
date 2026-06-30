import { useState, useEffect } from "react";
import { getActiveLoans, returnBook, type Loan } from "../../lib/api";
import LoanTable from "./LoanTable";

export default function ReturnForm() {
  const [loans, setLoans] = useState<Loan[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [returningId, setReturningId] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const loadLoans = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await getActiveLoans();
      setLoans(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load active loans");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadLoans();
  }, []);

  const handleReturn = async (loan: Loan) => {
    setReturningId(loan.id);
    setError(null);
    setSuccess(null);
    try {
      await returnBook(loan.id);
      setSuccess(`Returned "${loan.book_title}" successfully`);
      setLoans((prev) => prev.filter((l) => l.id !== loan.id));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to process return");
    } finally {
      setReturningId(null);
    }
  };

  if (loading) {
    return (
      <div class="flex items-center justify-center py-8">
        <div class="h-6 w-6 animate-spin rounded-full border-2 border-book-cloth border-t-leather" />
        <span class="ml-2 text-reading-lamp text-body">Loading active loans...</span>
      </div>
    );
  }

  return (
    <div class="space-y-4">
      {error && <p class="text-red-700 text-body bg-red-50 rounded p-3">{error}</p>}
      {success && <p class="text-green-700 text-body bg-green-50 rounded p-3">{success}</p>}
      <LoanTable loans={loans} onReturn={handleReturn} returningId={returningId ?? undefined} />
    </div>
  );
}
