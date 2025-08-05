import { SignupForm } from "../components/SignUpForm";
import { LinkIcon } from "lucide-react";
import { Link } from "react-router-dom";

export const SignUpPage = () => {
  return (
    <div
      className="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8"
      style={{ backgroundColor: "#F2EFDE" }}
    >
      <div className="max-w-md w-full space-y-8">
        <div className="text-center">
          <div className="flex justify-center items-center mb-6">
            <div
              className="p-4 rounded-2xl shadow-lg"
              style={{
                background: `linear-gradient(135deg, #A4193D 0%, #F4A261 100%)`,
              }}
            >
              <LinkIcon className="h-12 w-12 text-white" />
            </div>
          </div>
          <h2
            className="text-4xl font-bold mb-2 hover:bg-amber-500"
            style={{
              background: `linear-gradient(135deg, #A4193D 0%, #F4A261 100%)`,
              WebkitBackgroundClip: "text",
              WebkitTextFillColor: "transparent",
              backgroundClip: "text",
            }}
          >
            Create Account
          </h2>
          <p className="text-lg" style={{ color: "#6B7280" }}>
            Start shortening URLs and tracking analytics
          </p>
        </div>

        <div className="bg-white rounded-2xl shadow-xl p-8 border border-gray-100">
          <SignupForm />

          <div className="relative my-6">
            <div className="absolute inset-0 flex items-center">
              <div className="w-full border-t border-gray-200" />
            </div>
            <div className="relative flex justify-center text-sm">
              <span className="px-4 bg-white text-gray-500 font-medium">
                Already have an account?
              </span>
            </div>
          </div>

          <div className="text-center">
            <Link
              to="/login"
              className="inline-flex items-center px-6 py-2 border-2 rounded-xl font-semibold text-sm transition-all duration-200 hover:shadow-md border-[#A4193D] text-[#A4193D] hover:bg-[#A4193D] hover:text-white"
            >
              Sign In Instead
            </Link>
          </div>
        </div>

        <div className="text-center">
          <p className="text-sm" style={{ color: "#6B7280" }}>
            Secure signup powered by Supabase
          </p>
        </div>
      </div>
    </div>
  );
};
