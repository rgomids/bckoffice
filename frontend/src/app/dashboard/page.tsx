import ProtectedRoute from "@/components/ProtectedRoute";

export default function DashboardPage() {
  return (
    <ProtectedRoute>
      <div className="p-4">Dashboard</div>
    </ProtectedRoute>
  );
}
