import { useState, useEffect } from "react";
import { listPatrons } from "../../lib/api";

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

interface FineListProps {
  fines: Fine[];
  onPay?: (fineId: string) => void;
  payingId?: string | null;
}

export default function FineList({ fines, onPay, payingId }: FineListProps) {
  if (fines.length === 0) {
    return <p class="text-reading-lamp text-body">No fines found.</p>;
  }

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat("en-US", { style: "currency", currency: "USD" }).format(amount);
  };

  return (
    <div class="overflow-x-auto">
      <table class="w-full text-left text-body">
        <thead>
          <tr class="border-b border-book-cloth text-data text-reading-lamp">
            <th class="py-2 pr-4">Patron</th>
            <th class="py-2 pr-4">Book</th>
            <th class="py-2 pr-4">Amount</th>
            <th class="py-2 pr-4">Status</th>
            <th class="py-2 pr-4">Due Date</th>
            {onPay && <th class="py-2">Action</th>}
          </tr>
        </thead>
        <tbody>
          {fines.map((fine) => (
            <tr key={fine.id} class="border-b border-book-cloth/50">
              <td class="py-3 pr-4 font-medium text-ink">{fine.patron_name}</td>
              <td class="py-3 pr-4 text-reading-lamp">{fine.book_title}</td>
              <td class="py-3 pr-4 font-mono text-data">{formatCurrency(fine.amount)}</td>
              <td class="py-3 pr-4">
                <span class={["text-small px-2 py-0.5 rounded", fine.paid ? "bg-green-900/20 text-green-700" : "bg-red-900/20 text-red-700"].join(" ")}>
                  {fine.paid ? "Paid" : "Unpaid"}
                </span>
              </td>
              <td class="py-3 pr-4 text-reading-lamp">{new Date(fine.due_date).toLocaleDateString()}</td>
              {onPay && !fine.paid && (
                <td class="py-3">
                  <button
                    onClick={() => onPay(fine.id)}
                    disabled={payingId === fine.id}
                    class="px-3 py-1 rounded bg-leather text-paper text-small hover:bg-leather/90 transition-colors disabled:opacity-50"
                  >
                    {payingId === fine.id ? "Processing..." : "Mark Paid"}
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
