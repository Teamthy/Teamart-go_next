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

      <header className="flex min-h-screen flex-col items-center overflow-hidden bg-gradient-to-b from-slate-950 via-slate-900 to-slate-950 text-white">
        <nav className="sticky top-0 z-50 w-full border-b border-white/10 bg-white/5 backdrop-blur-md">
          <div className="flex items-center justify-between px-6 py-4 md:px-12 lg:px-20">
            <Link href="/" className="flex items-center gap-3">
              <div className="flex size-10 items-center justify-center rounded-xl bg-indigo-600 text-lg font-bold shadow-lg shadow-indigo-600/30">
                T
              </div>

              <div>
                <h1 className="text-lg font-semibold tracking-tight">Teamart</h1>
                <p className="-mt-1 text-xs text-slate-400">AI Commerce Platform</p>
              </div>
            </Link>

            <div
              className={`
                ${mobileOpen ? "max-md:w-full" : "max-md:w-0"}
                flex items-center gap-8 text-sm
                max-md:fixed
                max-md:left-0
                max-md:top-0
                max-md:h-screen
                max-md:flex-col
                max-md:justify-center
                max-md:overflow-hidden
                max-md:bg-slate-950/95
                max-md:backdrop-blur-xl
                max-md:transition-all
                max-md:duration-300
              `}
            >
              <Link href="#features" onClick={() => setMobileOpen(false)} className="text-slate-300 transition hover:text-white">
                Features
              </Link>
              <Link href="#creators" onClick={() => setMobileOpen(false)} className="text-slate-300 transition hover:text-white">
                Creators
              </Link>
              <Link href="#commerce" onClick={() => setMobileOpen(false)} className="text-slate-300 transition hover:text-white">
                Commerce
              </Link>
              <Link href="#pricing" onClick={() => setMobileOpen(false)} className="text-slate-300 transition hover:text-white">
                Pricing
              </Link>

              <button onClick={() => setMobileOpen(false)} className="rounded-lg bg-white p-2 text-black md:hidden">
                ✕
              </button>
            </div>

            <div className="hidden items-center gap-4 md:flex">
              <Link href="/dashboard" className="rounded-lg border border-white/10 px-5 py-2.5 transition hover:bg-white/5">
                Dashboard
              </Link>
              <Link href="/products" className="rounded-lg bg-indigo-600 px-5 py-2.5 font-medium shadow-lg shadow-indigo-600/20 transition hover:bg-indigo-500">
                Launch Store
              </Link>
            </div>

            <button onClick={() => setMobileOpen(true)} className="rounded-lg border border-white/10 bg-white/10 p-2 md:hidden">
              ☰
            </button>
          </div>
        </nav>

        <section className="mx-auto grid w-full max-w-7xl items-center gap-16 px-6 py-20 lg:grid-cols-2 lg:px-20 lg:py-28">
          <div>
            <div className="inline-flex items-center gap-2 rounded-full border border-indigo-500/20 bg-indigo-500/10 px-4 py-2 text-sm text-indigo-300">
              <span className="size-2 animate-pulse rounded-full bg-indigo-400"></span>
              AI Native Commerce Infrastructure
            </div>

            <h1 className="mt-6 text-5xl font-semibold leading-tight tracking-tight md:text-6xl">
              Build the next generation of <span className="text-indigo-400">creator commerce</span> with AI automation
            </h1>

            <p className="mt-6 max-w-2xl text-lg leading-8 text-slate-400">
              Teamart helps creators, brands, and digital businesses launch scalable AI-powered storefronts, automate workflows,
              manage users, process payments, and grow revenue from one unified platform.
            </p>

            <div className="mt-8 flex flex-wrap gap-4">
              <Link href="/auth/register" className="inline-block rounded-xl bg-indigo-600 px-7 py-4 font-medium shadow-xl shadow-indigo-600/20 transition hover:bg-indigo-500">
                Get Started
              </Link>
              <Link href="/docs" className="inline-block rounded-xl border border-white/10 px-7 py-4 transition hover:bg-white/5">
                View Documentation
              </Link>
            </div>

            <div className="mt-12 grid max-w-xl grid-cols-3 gap-6">
              <div>
                <h2 className="text-3xl font-bold">10K+</h2>
                <p className="mt-1 text-sm text-slate-400">Active Users</p>
              </div>
              <div>
                <h2 className="text-3xl font-bold">500+</h2>
                <p className="mt-1 text-sm text-slate-400">Creator Stores</p>
              </div>
              <div>
                <h2 className="text-3xl font-bold">$2M+</h2>
                <p className="mt-1 text-sm text-slate-400">Revenue Processed</p>
              </div>
            </div>
          </div>

          <div className="relative">
            <div className="absolute inset-0 rounded-full bg-indigo-600/20 blur-3xl"></div>
            <div className="relative rounded-3xl border border-white/10 bg-white/5 p-6 shadow-2xl backdrop-blur-xl">
              <div className="flex items-center justify-between border-b border-white/10 pb-4">
                <div>
                  <h3 className="text-lg font-semibold">Analytics Dashboard</h3>
                  <p className="text-sm text-slate-400">Real-time commerce insights</p>
                </div>
                <div className="rounded-full bg-emerald-500/20 px-3 py-1 text-sm text-emerald-400">Live</div>
              </div>

              <div className="mt-6 grid grid-cols-2 gap-4">
                <div className="rounded-2xl border border-white/5 bg-slate-900/70 p-5">
                  <p className="text-sm text-slate-400">Monthly Revenue</p>
                  <h2 className="mt-2 text-3xl font-bold">$48,320</h2>
                  <p className="mt-2 text-sm text-emerald-400">+18.4% growth</p>
                </div>

                <div className="rounded-2xl border border-white/5 bg-slate-900/70 p-5">
                  <p className="text-sm text-slate-400">AI Orders</p>
                  <h2 className="mt-2 text-3xl font-bold">12,847</h2>
                  <p className="mt-2 text-sm text-indigo-400">Automated workflows</p>
                </div>

                <div className="col-span-2 rounded-2xl border border-white/5 bg-slate-900/70 p-5">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="text-sm text-slate-400">User Growth</p>
                      <h2 className="mt-2 text-3xl font-bold">+324%</h2>
                    </div>
                    <div className="flex h-24 w-40 items-end gap-2 rounded-xl border border-white/10 bg-gradient-to-tr from-indigo-500/20 to-cyan-500/20 p-3">
                      <div className="h-10 w-4 rounded bg-indigo-400"></div>
                      <div className="h-16 w-4 rounded bg-indigo-400"></div>
                      <div className="h-12 w-4 rounded bg-indigo-400"></div>
                      <div className="h-20 w-4 rounded bg-indigo-400"></div>
                      <div className="h-14 w-4 rounded bg-indigo-400"></div>
                    </div>
                  </div>
                </div>
              </div>

              <div className="mt-6 rounded-2xl border border-white/5 bg-slate-900/70 p-5">
                <div className="flex items-center justify-between">
                  <div>
                    <h3 className="font-semibold">Registered Users</h3>
                    <p className="text-sm text-slate-400">Connected accounts across the platform</p>
                  </div>
                  <button className="text-sm text-indigo-400 hover:text-indigo-300">View All</button>
                </div>

                <div className="mt-4 space-y-3">
                  {[
                    "Creator Economy",
                    "Digital Products",
                    "Affiliate Commerce",
                    "AI Automation",
                  ].map((item, index) => (
                    <div key={index} className="flex items-center justify-between rounded-xl bg-white/5 px-4 py-3">
                      <div className="flex items-center gap-3">
                        <div className="flex size-10 items-center justify-center rounded-full bg-indigo-500/20">🚀</div>
                        <div>
                          <h4 className="font-medium">{item}</h4>
                          <p className="text-xs text-slate-400">Active Module</p>
                        </div>
                      </div>
                      <span className="text-sm text-emerald-400">Online</span>
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
