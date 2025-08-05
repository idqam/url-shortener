import { useAuthStore } from "../store/AuthStore";
import { Navigate } from "react-router-dom";

interface ProtectedRouteProps {
  children: React.ReactNode;
}

export const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children }) => {
  const { isAuthenticated, accessToken } = useAuthStore();

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  if (!accessToken) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <div className="flex space-x-2 justify-center m-0 items-center w-fit h-fit ">
            <p>
              Loading
              <span className="h-2 w-2 bg-amber-300 rounded-full animate-pulse">
                .
              </span>
              <span className="h-2 w-2 bg-amber-300 rounded-full animate-pulse">
                .
              </span>
              <span className="h-2 w-2 bg-amber-300 rounded-full animate-pulse">
                .
              </span>
            </p>
          </div>
        </div>
      </div>
    );
  }

  return <>{children}</>;
};
