import { SignupForm } from "../components/SignUpForm";
import { Link, useNavigate } from "react-router-dom";
import { Link as LinkIcon } from "lucide-react";

export const SignUpPage = () => {
  const navigate = useNavigate();

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col">
      <header className="bg-white border-b border-gray-200">
        <div className="max-w-md mx-auto px-6 h-14 flex items-center">
          <div className="flex items-center gap-2 cursor-pointer" onClick={() => navigate("/")}>
            <div className="w-7 h-7 bg-blue-600 rounded-md flex items-center justify-center">
              <LinkIcon className="w-4 h-4 text-white" />
            </div>
            <span className="font-semibold text-gray-900 text-sm">Shortlink</span>
          </div>
        </div>
      </header>

      <main className="flex-1 flex items-center justify-center px-4 py-12">
        <div className="max-w-md w-full">
          <div className="mb-8">
            <h1 className="text-2xl font-bold text-gray-900">Create an account</h1>
            <p className="text-sm text-gray-500 mt-1">Start shortening URLs and tracking analytics</p>
          </div>

          <div className="bg-white rounded-lg border border-gray-200 p-8">
            <SignupForm />

            <p className="mt-6 text-center text-sm text-gray-500">
              Already have an account?{" "}
              <Link to="/login" className="font-medium text-blue-600 hover:text-blue-700">
                Sign in
              </Link>
            </p>
          </div>
        </div>
      </main>
    </div>
  );
};
