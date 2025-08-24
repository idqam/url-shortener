import { useNavigate } from "react-router-dom";
import { ShortenerForm } from "../components/ShortenerForm";

export const Home = () => {
  const navigate = useNavigate();


  return (
    <div
      className="min-h-screen relative overflow-hidden"
      style={{ backgroundColor: "#F2EFDE" }}
    >
    

      <section className="relative z-10 px-6 py-8">
        <div className="flex justify-center pt-8 mb-16">
          <div className="relative group">
            <div className="bg-white/70 backdrop-blur-lg rounded-3xl shadow-2xl border border-white/20 p-8 hover:shadow-3xl transition-all duration-500 hover:scale-105">
              <div
                className="absolute inset-0 rounded-3xl opacity-10"
                style={{
                  background: "linear-gradient(135deg, #A4193D, #F4A261)",
                }}
              ></div>
              <div className="relative">
                <div className="flex items-center justify-center space-x-4 mb-4">
                  <div
                    className="w-12 h-12 rounded-2xl flex items-center justify-center shadow-lg group-hover:rotate-12 transition-transform duration-300"
                    style={{
                      background: "linear-gradient(135deg, #A4193D, #F4A261)",
                    }}
                  >
                    <svg
                      className="w-6 h-6 text-white"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"
                      />
                    </svg>
                  </div>
                  <div
                    className="text-4xl font-black bg-clip-text text-transparent"
                    style={{
                      backgroundImage:
                        "linear-gradient(135deg, #A4193D, #F4A261)",
                    }}
                  >
                    URL SHORTENER
                  </div>
                </div>
                <div
                  style={{ color: "#6B7280" }}
                  className="text-center text-lg font-medium"
                >
                  Transform long URLs into powerful short links
                </div>
              </div>
            </div>
          </div>
        </div>

        <div className="max-w-6xl mx-auto">
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 mb-12 items-start">
            <div className="lg:col-span-2">
              <div className="bg-white/70 backdrop-blur-lg rounded-3xl shadow-xl border border-white/20 p-8 hover:shadow-2xl transition-all duration-500">
                <div
                  className="absolute inset-0 rounded-3xl opacity-5"
                  style={{
                    background: "linear-gradient(135deg, #A4193D, #F4A261)",
                  }}
                ></div>
                <div className="relative">
                  <h2
                    className="text-2xl font-bold mb-6 text-center"
                    style={{ color: "#111827" }}
                  >
                    Shorten Your URL
                  </h2>
                  <ShortenerForm />
                </div>
              </div>
            </div>

            <div
              className="rounded-3xl shadow-xl p-8 text-white relative overflow-hidden hover:scale-105 transition-all duration-500 h-fit"
              style={{
                background: "linear-gradient(135deg, #A4193D, #F4A261)",
              }}
            >
              <div className="absolute top-0 right-0 w-32 h-32 bg-white/10 rounded-full -translate-y-16 translate-x-16"></div>
              <div className="absolute bottom-0 left-0 w-24 h-24 bg-white/10 rounded-full translate-y-12 -translate-x-12"></div>
              <div className="relative z-10">
                <div className="w-12 h-12 bg-white/20 rounded-2xl flex items-center justify-center mb-4">
                  <svg
                    className="w-6 h-6"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M13 10V3L4 14h7v7l9-11h-7z"
                    />
                  </svg>
                </div>
                <h3 className="text-xl font-bold mb-2">Lightning Fast</h3>
                <p className="text-white/90 text-sm">
                  Generate short links quickly and see basic core analytics in
                  your dashboard by making an account!
                </p>
              </div>
            </div>
          </div>

          <div className="max-w-4xl mx-auto">
            <div className="bg-white/70 backdrop-blur-lg rounded-3xl shadow-xl border border-white/20 p-8 text-center relative overflow-hidden">
              <div
                className="absolute inset-0 rounded-3xl opacity-5"
                style={{
                  background: "linear-gradient(135deg, #A4193D, #F4A261)",
                }}
              ></div>
              <div className="relative">
                <div className="mb-6">
                  <h2
                    className="text-3xl font-bold mb-4"
                    style={{ color: "#A4193D" }}
                  >
                    Premium features
                  </h2>

                  <div
                    className="inline-block px-4 py-2 rounded-full text-sm font-semibold text-white mb-4"
                    style={{
                      background: "linear-gradient(135deg, #A4193D, #F4A261)",
                    }}
                  >
                    Coming Soon
                  </div>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
                  {[
                    { title: "Custom Domains", desc: "brand.ly/link" },
                    {
                      title: "Advanced Analytics",
                      desc: "Fine-grain insights and trends",
                    },
                    { title: "QR Code Generator", desc: "Instant QR codes" },
                  ].map((item, index) => (
                    <div
                      key={index}
                      className="bg-white/50 rounded-xl p-4 border border-white/30"
                    >
                      <h4
                        className="font-semibold mb-1"
                        style={{ color: "#A4193D" }}
                      >
                        {item.title}
                      </h4>
                      <p className="text-sm" style={{ color: "#6B7280" }}>
                        {item.desc}
                      </p>
                    </div>
                  ))}
                </div>

                <div className="flex flex-wrap justify-center items-center gap-4">
                  <button
                    className="font-semibold py-3 px-8 rounded-xl transition-all duration-300 hover:scale-105 hover:shadow-lg text-white"
                    style={{
                      background: "linear-gradient(135deg, #A4193D, #F4A261)",
                    }}
                    onClick={() => navigate("/signup")}
                  >
                    Create Account
                  </button>
                  <button
                    className="font-semibold py-3 px-8 rounded-xl transition-all duration-300 hover:scale-105 border-2 hover:shadow-lg"
                    style={{
                      borderColor: "#A4193D",
                      color: "#A4193D",
                      backgroundColor: "transparent",
                    }}
                    onMouseEnter={(e) => {
                      const target = e.target as HTMLButtonElement;
                      target.style.backgroundColor = "#A4193D";
                      target.style.color = "white";
                    }}
                    onMouseLeave={(e) => {
                      const target = e.target as HTMLButtonElement;
                      target.style.backgroundColor = "transparent";
                      target.style.color = "#A4193D";
                    }}
                    onClick={() => navigate("/login")}
                  >
                    Sign In
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div className="flex justify-center mt-16">
          <div className="max-w-4xl w-full bg-white/50 backdrop-blur-sm rounded-2xl border border-white/30 p-6 text-center">
            <div className="flex items-center justify-center mb-3">
              <h3
                className="text-lg font-semibold"
                style={{ color: "#A4193D" }}
              >
                Privacy & Analytics Notice
              </h3>
            </div>
            <p className="text-sm mb-4" style={{ color: "#6B7280" }}>
              We collect basic analytics to improve our service: click counts,
              referrer sources, device types, and daily usage trends. No
              personal information is stored with your shortened URLs.
            </p>
            <div
              className="flex flex-wrap justify-center gap-4 text-xs"
              style={{ color: "#6B7280" }}
            >
              <span className="flex items-center">
                <div
                  className="w-2 h-2 rounded-full mr-2"
                  style={{ backgroundColor: "#2A9D8F" }}
                ></div>
                Click Analytics
              </span>
              <span className="flex items-center">
                <div
                  className="w-2 h-2 rounded-full mr-2"
                  style={{ backgroundColor: "#F4A261" }}
                ></div>
                Device Stats
              </span>
              <span className="flex items-center">
                <div
                  className="w-2 h-2 rounded-full mr-2"
                  style={{ backgroundColor: "#A4193D" }}
                ></div>
                Referrer Tracking
              </span>
              <span className="flex items-center">
                <div
                  className="w-2 h-2 rounded-full mr-2"
                  style={{ backgroundColor: "#FFDFB9" }}
                ></div>
                Usage Trends
              </span>
            </div>
          </div>
        </div>
      </section>
    </div>
  );
};
