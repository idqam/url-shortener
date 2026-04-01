import { useNavigate } from "react-router-dom";
import { ShortenerForm } from "../components/ShortenerForm";
import { Link as LinkIcon } from "lucide-react";

export const Home = () => {
  const navigate = useNavigate();

  return (
    <div className="min-h-screen bg-white flex flex-col">
      {/* Nav */}
      <header className="border-b border-gray-200">
        <div className="max-w-3xl mx-auto px-6 h-14 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="w-7 h-7 bg-blue-600 rounded-md flex items-center justify-center">
              <LinkIcon className="w-4 h-4 text-white" />
            </div>
            <span className="font-semibold text-gray-900 text-sm">Shortlink</span>
          </div>
          <div className="flex items-center gap-3">
            <button
              onClick={() => navigate("/login")}
              className="text-sm font-medium text-gray-600 hover:text-gray-900 transition-colors"
            >
              Sign in
            </button>
            <button
              onClick={() => navigate("/signup")}
              className="text-sm font-medium bg-blue-600 text-white px-3 py-1.5 rounded-md hover:bg-blue-700 transition-colors"
            >
              Get started
            </button>
          </div>
        </div>
      </header>

      {/* Hero */}
      <main className="flex-1 flex flex-col items-center px-6 pt-20 pb-16">
        <div className="max-w-xl w-full text-center mb-10">
          <h1 className="text-3xl font-bold text-gray-900 mb-3">Shorten any URL</h1>
          <p className="text-gray-500 text-base">
            Paste a link, get a short one. Sign up to track clicks with analytics.
          </p>
        </div>

        <div className="max-w-xl w-full">
          <ShortenerForm />
        </div>

        {/* Upcoming features */}
        <div className="max-w-xl w-full mt-20 border-t border-gray-100 pt-10">
          <p className="text-xs font-medium text-gray-400 uppercase tracking-wider mb-6 text-center">
            Coming soon
          </p>
          <div className="grid grid-cols-3 gap-6">
            {[
              { title: "Custom Domains", desc: "Use your own domain" },
              { title: "Advanced Analytics", desc: "Fine-grain click insights" },
              { title: "QR Codes", desc: "Instant QR generation" },
            ].map((item) => (
              <div key={item.title} className="text-center">
                <p className="text-sm font-medium text-gray-700">{item.title}</p>
                <p className="text-xs text-gray-400 mt-0.5">{item.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </main>

      {/* Footer */}
      <footer className="border-t border-gray-100 py-4">
        <p className="text-xs text-gray-400 text-center">
          We collect basic analytics: click counts, referrer sources, device types. No personal data stored with your links.
        </p>
      </footer>
    </div>
  );
};
