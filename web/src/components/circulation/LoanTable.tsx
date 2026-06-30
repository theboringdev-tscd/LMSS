import { useState, useEffect } from "react";
import { getActiveLoans, type Loan } from "../../lib/api";

interface LoanTableProps {
  loans: Loan[];
  onReturn?: (loan: Loan) => void;
  returningId?: string;
}

export default function LoanTable({ loans, onReturn, returningId }: LoanTableProps) {
  if (loans.length === 0) {
    return <p class="text-reading-lamp text-body">No loans found.</p>;
  }

  return (
    <div class="overflow-x-auto">
      <table class="w-full text-left text-body">
        <thead>
          <tr class="border-b border-book-cloth text-data text-reading-lamp">
            <th class="py-2 pr-4">Book</th>
            <th class="py-2 pr-4">Patron</th>
            <th class="py-2 pr-4">Checked Out</th>
            <th class="py-2 pr-4">Due Date</th>
            <th class="py-2 pr-4">Status</th>
            {onReturn && <th class="py-2">Action</th>}
          </tr>
        </thead>
        <tbody>
          {loans.map((loan) => (
            <tr key={loan.id} class="border-b border-book-cloth/50">
              <td class="py-3 pr-4 font-medium text-ink">{loan.book_title}</td>
              <td class="py-3 pr-4 text-reading-lamp">{loan.patron_name}</td>
              <td class="py-3 pr-4 text-reading-lamp">{new Date(loan.checked_out_at).toLocaleDateString()}</td>
              <td class="py-3 pr-4 font-mono text-data text-reading-lamp">{new Date(loan.due_date).toLocaleDateString()}</td>
              <td class="py-3 pr-4">
                <span
                  class={[
                    "text-small px-2 py-0.5 rounded",
                    loan.status === "active" ? "bg-leather/20 text-leather" : "bg-green-900/20 text-green-700",
                  ].join(" ")}
                >
                  {loan.status}
                </span>
              </td>
              {onReturn && (
                <td class="py-3">
                  <button
                    onClick={() => onReturn(loan)}
                    disabled={returningId === loan.id}
                    class="px-3 py-1 rounded bg-book-cloth text-ink text-small hover:bg-book-cloth/80 transition-colors disabled:opacity-50"
                  >
                    {returningId === loan.id ? "Processing..." : "Return"}
                  </button>
                </td>
              )}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
