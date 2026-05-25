"use client";

import React from "react";
import Link from "next/link";

export default function Home() {
  const [mobileOpen, setMobileOpen] = React.useState(false);

  return (
    <>
      <style>
        {`
          @import url("https://fonts.googleapis.com/css2?family=Poppins:wght@300;400;500;600;700&display=swap");

          * {
            font-family: "Poppins", sans-serif;
          }

          body {
            background: #0b1120;
          }
        `}
      </style>

      {/* HERO SECTION */}
      <header className="flex flex-col items-center bg-gradient-to-b from-slate-950 via-slate-900 to-slate-950 text-white min-h-screen overflow-hidden">

        {/* NAVBAR */}
        <nav className="w-full border-b border-white/10 backdrop-blur-md bg-white/5 sticky top-0 z-50">
          <div className="flex items-center justify-between px-6 md:px-12 lg:px-20 py-4">

            {/* LOGO */}
            <Link href="/" className="flex items-center gap-3">
              <div className="size-10 rounded-xl bg-indigo-600 flex items-center justify-center text-lg font-bold shadow-lg shadow-indigo-600/30">
                T
              </div>

              <div>
                <h1 className="text-lg font-semibold tracking-tight">
                  Teamart
                </h1>
                <p className="text-xs text-slate-400 -mt-1">
                  AI Commerce Platform
                </p>
              </div>
            </Link>

            {/* MENU */}
            <div
              className={`
                ${mobileOpen ? "max-md:w-full" : "max-md:w-0"}
                max-md:fixed
                max-md:top-0
                max-md:left-0
                max-md:h-screen
                max-md:bg-slate-950/95
                max-md:backdrop-blur-xl
                max-md:flex-col
                max-md:justify-center
                max-md:overflow-hidden
                max-md:transition-all
                max-md:duration-300
                flex items-center gap-8 text-sm
              `}
            >
              <Link
                href="#features"
                onClick={() => setMobileOpen(false)}
                className="text-slate-300 hover:text-white transition"
              >
                Features
              </Link>

              <Link
                href="#creators"
                onClick={() => setMobileOpen(false)}
                className="text-slate-300 hover:text-white transition"
              >
                Creators
              </Link>

              <Link
                href="#commerce"
                onClick={() => setMobileOpen(false)}
                className="text-slate-300 hover:text-white transition"
              >
                Commerce
              </Link>

              <Link
                href="#pricing"
                onClick={() => setMobileOpen(false)}
                className="text-slate-300 hover:text-white transition"
              >
                Pricing
              </Link>

              <button
                onClick={() => setMobileOpen(false)}
                className="md:hidden bg-white text-black p-2 rounded-lg"
              >
                ✕
              </button>
            </div>

            {/* CTA */}
            <div className="hidden md:flex items-center gap-4">
              <Link href="/dashboard" className="px-5 py-2.5 rounded-lg border border-white/10 hover:bg-white/5 transition">
                Dashboard
              </Link>

              <Link href="/products" className="bg-indigo-600 hover:bg-indigo-500 transition px-5 py-2.5 rounded-lg shadow-lg shadow-indigo-600/20 font-medium">
                Launch Store
              </Link>
            </div>

            {/* MOBILE MENU BTN */}
            <button
              onClick={() => setMobileOpen(true)}
              className="md:hidden bg-white/10 border border-white/10 p-2 rounded-lg"
            >
              ☰
            </button>
          </div>
        </nav>

        {/* HERO CONTENT */}
        <section className="w-full max-w-7xl mx-auto px-6 md:px-12 lg:px-20 py-20 lg:py-28 grid lg:grid-cols-2 gap-16 items-center">

          {/* LEFT */}
          <div>

            {/* BADGE */}
            <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-indigo-500/10 border border-indigo-500/20 text-indigo-300 text-sm">
              <span className="size-2 rounded-full bg-indigo-400 animate-pulse"></span>
              AI Native Commerce Infrastructure
            </div>

            {/* HEADING */}
            <h1 className="mt-6 text-5xl md:text-6xl font-semibold tracking-tight leading-tight">
              Build the next generation of{" "}
              <span className="text-indigo-400">
                creator commerce
              </span>{" "}
              with AI automation
            </h1>

            {/* DESCRIPTION */}
            <p className="mt-6 text-slate-400 text-lg leading-8 max-w-2xl">
              Teamart helps creators, brands, and digital businesses launch
              scalable AI-powered storefronts, automate workflows, manage users,
              process payments, and grow revenue from one unified platform.
            </p>

            {/* CTA BUTTONS */}
            <div className="mt-8 flex flex-wrap gap-4">
              <Link href="/auth/signup" className="bg-indigo-600 hover:bg-indigo-500 transition px-7 py-4 rounded-xl font-medium shadow-xl shadow-indigo-600/20 inline-block">
                Get Started
              </Link>

              <Link href="/docs" className="border border-white/10 hover:bg-white/5 transition px-7 py-4 rounded-xl inline-block">
                View Documentation
              </Link>
            </div>

            {/* STATS */}
            <div className="grid grid-cols-3 gap-6 mt-12 max-w-xl">

              <div>
                <h2 className="text-3xl font-bold">10K+</h2>
                <p className="text-slate-400 text-sm mt-1">
                  Active Users
                </p>
              </div>

              <div>
                <h2 className="text-3xl font-bold">500+</h2>
                <p className="text-slate-400 text-sm mt-1">
                  Creator Stores
                </p>
              </div>

              <div>
                <h2 className="text-3xl font-bold">$2M+</h2>
                <p className="text-slate-400 text-sm mt-1">
                  Revenue Processed
                </p>
              </div>
            </div>
          </div>

          {/* RIGHT SIDE */}
          <div className="relative">

            {/* GLOW */}
            <div className="absolute inset-0 bg-indigo-600/20 blur-3xl rounded-full"></div>

            {/* DASHBOARD CARD */}
            <div className="relative bg-white/5 border border-white/10 backdrop-blur-xl rounded-3xl p-6 shadow-2xl">

              {/* TOP BAR */}
              <div className="flex items-center justify-between border-b border-white/10 pb-4">
                <div>
                  <h3 className="font-semibold text-lg">
                    Analytics Dashboard
                  </h3>
                  <p className="text-slate-400 text-sm">
                    Real-time commerce insights
                  </p>
                </div>

                <div className="bg-emerald-500/20 text-emerald-400 px-3 py-1 rounded-full text-sm">
                  Live
                </div>
              </div>

              {/* METRICS */}
              <div className="grid grid-cols-2 gap-4 mt-6">

                <div className="bg-slate-900/70 rounded-2xl p-5 border border-white/5">
                  <p className="text-slate-400 text-sm">
                    Monthly Revenue
                  </p>

                  <h2 className="text-3xl font-bold mt-2">
                    $48,320
                  </h2>

                  <p className="text-emerald-400 text-sm mt-2">
                    +18.4% growth
                  </p>
                </div>

                <div className="bg-slate-900/70 rounded-2xl p-5 border border-white/5">
                  <p className="text-slate-400 text-sm">
                    AI Orders
                  </p>

                  <h2 className="text-3xl font-bold mt-2">
                    12,847
                  </h2>

                  <p className="text-indigo-400 text-sm mt-2">
                    Automated workflows
                  </p>
                </div>

                <div className="bg-slate-900/70 rounded-2xl p-5 border border-white/5 col-span-2">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="text-slate-400 text-sm">
                        User Growth
                      </p>

                      <h2 className="text-3xl font-bold mt-2">
                        +324%
                      </h2>
                    </div>

                    <div className="h-24 w-40 rounded-xl bg-gradient-to-tr from-indigo-500/20 to-cyan-500/20 border border-white/10 flex items-end gap-2 p-3">
                      <div className="w-4 h-10 bg-indigo-400 rounded"></div>
                      <div className="w-4 h-16 bg-indigo-400 rounded"></div>
                      <div className="w-4 h-12 bg-indigo-400 rounded"></div>
                      <div className="w-4 h-20 bg-indigo-400 rounded"></div>
                      <div className="w-4 h-14 bg-indigo-400 rounded"></div>
                    </div>
                  </div>
                </div>
              </div>

              {/* BOTTOM USERS */}
              <div className="mt-6 bg-slate-900/70 rounded-2xl border border-white/5 p-5">
                <div className="flex items-center justify-between">
                  <div>
                    <h3 className="font-semibold">
                      Registered Users
                    </h3>

                    <p className="text-slate-400 text-sm">
                      Connected accounts across the platform
                    </p>
                  </div>

                  <button className="text-indigo-400 text-sm hover:text-indigo-300">
                    View All
                  </button>
                </div>

                <div className="mt-4 space-y-3">

                  {[
                    "Creator Economy",
                    "Digital Products",
                    "Affiliate Commerce",
                    "AI Automation",
                  ].map((item, index) => (
                    <div
                      key={index}
                      className="flex items-center justify-between bg-white/5 rounded-xl px-4 py-3"
                    >
                      <div className="flex items-center gap-3">
                        <div className="size-10 rounded-full bg-indigo-500/20 flex items-center justify-center">
                          🚀
                        </div>

                        <div>
                          <h4 className="font-medium">{item}</h4>
                          <p className="text-xs text-slate-400">
                            Active Module
                          </p>
                        </div>
                      </div>

                      <span className="text-emerald-400 text-sm">
                        Online
                      </span>
                    </div>
                  ))}

                </div>
              </div>
            </div>
          </div>
        </section>
      </header>
    </>
  );
}
            Deploy Now
          </a >
  <a
    className="flex h-12 w-full items-center justify-center rounded-full border border-solid border-black/[.08] px-5 transition-colors hover:border-transparent hover:bg-black/[.04] dark:border-white/[.145] dark:hover:bg-[#1a1a1a] md:w-[158px]"
    href="https://nextjs.org/docs?utm_source=create-next-app&utm_medium=appdir-template-tw&utm_campaign=create-next-app"
    target="_blank"
    rel="noopener noreferrer"
  >
    Documentation
  </a>
        </div >
      </main >
    </div >
  );
}
