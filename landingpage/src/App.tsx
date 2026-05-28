import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import {
  ChevronDown,
  ShoppingBag,
  Gift,
  Wallet,
  Shield,
  Clock,
  Car,
  Send,
  Receipt,
  Award,
  Briefcase,
  Globe,
  ArrowUp,
  ArrowRight,
  ArrowLeft,
  Play,
  Apple,
  QrCode,
  Menu,
  X,
  Bell,
  RefreshCw,
  CreditCard,
  Store,
} from 'lucide-react';

// ============ NAVIGATION ============
const products = [
  {
    icon: <ShoppingBag className="w-8 h-8 text-purple-600" />,
    iconBg: 'bg-purple-100',
    title: 'Blayzz Business',
    description: 'Connect your business with blayzz',
  },
  {
    icon: <Gift className="w-8 h-8 text-yellow-600" />,
    iconBg: 'bg-yellow-100',
    title: 'Refer & earn',
    description: 'Start earning by referring friends and family',
  },
  {
    icon: <Wallet className="w-8 h-8 text-pink-600" />,
    iconBg: 'bg-pink-100',
    title: 'Payment',
    description: 'Make free transfers, pay bills easily at no cost and pay in instalments',
  },
  {
    icon: <Shield className="w-8 h-8 text-teal-600" />,
    iconBg: 'bg-teal-100',
    title: 'Insurance',
    description: 'Buy insurance, view your certificates and make claims',
  },
  {
    icon: <Clock className="w-8 h-8 text-green-600" />,
    iconBg: 'bg-green-100',
    title: 'Budget & save',
    description: 'Create your savings plans and budget',
  },
  {
    icon: <Car className="w-8 h-8 text-blue-600" />,
    iconBg: 'bg-blue-100',
    title: 'Buy a car',
    description: 'Buy a car today. You can pay in instalments or pay outrightly',
  },
];

function Navbar() {
  const [showProducts, setShowProducts] = useState(false);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  return (
    <nav className="fixed top-0 left-0 right-0 z-50 px-4 py-3">
      <div className="max-w-6xl mx-auto bg-[#2d1010]/95 backdrop-blur-md rounded-xl border border-white/10 px-4 py-3 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <div className="flex items-center gap-1">
            <div className="relative">
              <div className="w-8 h-8 bg-red-600 rounded-lg flex items-center justify-center">
                <ArrowUp className="w-5 h-5 text-white" strokeWidth={3} />
              </div>
            </div>
            <span className="text-white font-bold text-xl ml-1">blayzz</span>
          </div>
          <span className="text-[10px] text-gray-400 hidden sm:block">By PremiumTrust Bank</span>
        </div>

        <div className="hidden md:flex items-center gap-6">
          <div className="relative">
            <button
              onMouseEnter={() => setShowProducts(true)}
              onMouseLeave={() => setShowProducts(false)}
              className="text-white hover:text-gray-300 flex items-center gap-1 text-sm font-medium"
            >
              Products <ChevronDown className="w-4 h-4" />
            </button>
            <AnimatePresence>
              {showProducts && (
                <motion.div
                  initial={{ opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: 10 }}
                  onMouseEnter={() => setShowProducts(true)}
                  onMouseLeave={() => setShowProducts(false)}
                  className="absolute top-full right-0 mt-2 w-[600px] bg-white rounded-xl shadow-2xl p-6 grid grid-cols-2 gap-4"
                >
                  {products.map((product, i) => (
                    <div key={i} className="flex items-start gap-3 p-3 rounded-lg hover:bg-gray-50 cursor-pointer transition-colors">
                      <div className={`w-12 h-12 ${product.iconBg} rounded-lg flex items-center justify-center flex-shrink-0`}>
                        {product.icon}
                      </div>
                      <div>
                        <h4 className="font-semibold text-gray-900 text-sm">{product.title}</h4>
                        <p className="text-xs text-gray-500 mt-0.5">{product.description}</p>
                      </div>
                    </div>
                  ))}
                </motion.div>
              )}
            </AnimatePresence>
          </div>
          <a href="#faq" className="text-white hover:text-gray-300 text-sm font-medium">FAQs</a>
          <button className="bg-red-600 hover:bg-red-700 text-white px-6 py-2 rounded-full text-sm font-medium transition-colors">
            Demo
          </button>
        </div>

        <button className="md:hidden text-white" onClick={() => setMobileMenuOpen(!mobileMenuOpen)}>
          {mobileMenuOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
        </button>
      </div>

      {/* Mobile Menu */}
      <AnimatePresence>
        {mobileMenuOpen && (
          <motion.div
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: 'auto' }}
            exit={{ opacity: 0, height: 0 }}
            className="md:hidden mt-2 bg-[#2d1010] rounded-xl p-4 space-y-4"
          >
            <a href="#" className="block text-white py-2">Products</a>
            <a href="#faq" className="block text-white py-2">FAQs</a>
            <button className="bg-red-600 text-white px-6 py-2 rounded-full w-full">Demo</button>
          </motion.div>
        )}
      </AnimatePresence>
    </nav>
  );
}

// ============ HERO SECTION ============
function HeroSection() {
  return (
    <section className="relative min-h-screen bg-gradient-to-b from-[#1a0000] via-[#8B0000] to-[#4A0000] pt-24 pb-20 overflow-hidden">
      {/* Background gradient overlay */}
      <div className="absolute inset-0 bg-gradient-to-r from-[#8B0000] via-[#C41E3A] to-[#8B0000] opacity-90" />
      
      <div className="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-12 lg:pt-20">
        <div className="grid lg:grid-cols-2 gap-12 items-center">
          {/* Left Content */}
          <div className="space-y-8">
            <h1 className="text-5xl sm:text-6xl lg:text-7xl xl:text-8xl font-bold text-white leading-[1.1] tracking-tight">
              Unlock
              <br />
              Financial
              <br />
              Freedom
            </h1>
            
            <div className="inline-flex items-center gap-3 bg-white rounded-full px-5 py-3 shadow-lg">
              <span className="text-gray-800 font-medium text-sm">Get the app</span>
              <Play className="w-4 h-4 text-green-600" fill="currentColor" />
              <Apple className="w-4 h-4 text-gray-800" fill="currentColor" />
              <QrCode className="w-4 h-4 text-gray-800" />
            </div>
          </div>

          {/* Right Content - Phone Mockup */}
          <div className="relative flex justify-center lg:justify-end">
            {/* Floating icons */}
            <motion.div
              animate={{ y: [0, -15, 0] }}
              transition={{ duration: 3, repeat: Infinity, ease: 'easeInOut' }}
              className="absolute -top-4 right-8 z-20 bg-white/10 backdrop-blur-sm rounded-2xl p-4 border border-white/20"
            >
              <div className="w-12 h-12 bg-red-600 rounded-xl flex items-center justify-center">
                <ArrowUp className="w-7 h-7 text-white" strokeWidth={3} />
              </div>
            </motion.div>

            <motion.div
              animate={{ y: [0, 10, 0] }}
              transition={{ duration: 4, repeat: Infinity, ease: 'easeInOut', delay: 0.5 }}
              className="absolute top-32 -left-4 z-20 bg-white/10 backdrop-blur-sm rounded-2xl p-4 border border-white/20"
            >
              <div className="w-14 h-14 bg-red-500/80 rounded-xl flex items-center justify-center">
                <ShoppingBag className="w-8 h-8 text-white" />
              </div>
            </motion.div>

            <motion.div
              animate={{ y: [0, -10, 0] }}
              transition={{ duration: 3.5, repeat: Infinity, ease: 'easeInOut', delay: 1 }}
              className="absolute bottom-32 -right-4 z-20 bg-white/10 backdrop-blur-sm rounded-2xl p-4 border border-white/20"
            >
              <div className="w-14 h-14 bg-orange-500/80 rounded-xl flex items-center justify-center">
                <Store className="w-8 h-8 text-white" />
              </div>
            </motion.div>

            {/* Phone Mockup */}
            <div className="relative z-10">
              <div className="w-[280px] sm:w-[320px] bg-gray-900 rounded-[3rem] p-3 shadow-2xl border-4 border-gray-800">
                <div className="bg-white rounded-[2.2rem] overflow-hidden">
                  {/* Phone Header */}
                  <div className="bg-white px-4 py-3 flex items-center justify-between border-b">
                    <div>
                      <p className="text-sm font-semibold text-gray-800">👋 Good morning</p>
                      <p className="text-[10px] text-gray-500">What would you like to do today</p>
                    </div>
                    <div className="flex gap-2">
                      <Bell className="w-5 h-5 text-yellow-500" />
                      <RefreshCw className="w-5 h-5 text-red-500" />
                    </div>
                  </div>

                  {/* Account Bar */}
                  <div className="bg-gray-100 px-4 py-2 flex items-center justify-between">
                    <span className="text-xs text-gray-600">My Account | 271652527</span>
                    <button className="text-xs text-gray-600 flex items-center gap-1">
                      View Balance <span className="text-[10px]">👁</span>
                    </button>
                  </div>

                  {/* Quick Actions */}
                  <div className="px-4 py-4 grid grid-cols-5 gap-2">
                    {[
                      { icon: '💳', label: 'Make Payment', badge: true },
                      { icon: '🎁', label: 'Refer & Get Reward', badge: true, badgeText: 'Freebies' },
                      { icon: '🛡️', label: 'Buy Insurance', badge: true, badgeText: '0%' },
                      { icon: '🔄', label: 'Renew Insurance' },
                      { icon: '📋', label: 'Make Claims' },
                    ].map((item, i) => (
                      <div key={i} className="flex flex-col items-center gap-1">
                        <div className="relative w-10 h-10 bg-gray-100 rounded-xl flex items-center justify-center text-lg">
                          {item.icon}
                          {item.badge && (
                            <span className="absolute -top-1 -right-1 w-4 h-4 bg-red-500 rounded-full text-[8px] text-white flex items-center justify-center">
                              {item.badgeText === 'Freebies' ? '🎉' : '0%'}
                            </span>
                          )}
                        </div>
                        <span className="text-[8px] text-gray-600 text-center leading-tight">{item.label}</span>
                      </div>
                    ))}
                  </div>

                  {/* Cards */}
                  <div className="px-4 pb-4 space-y-3">
                    <div className="grid grid-cols-2 gap-3">
                      <div className="bg-gray-800 rounded-xl p-3 text-white">
                        <h4 className="text-xs font-semibold mb-1">Transfer & Pay Bills</h4>
                        <p className="text-[9px] text-gray-300 mb-2">Unlimited free transfers. Pay bills at no cost.</p>
                        <div className="flex items-center justify-between">
                          <div className="flex items-center gap-1">
                            <span className="text-lg">📄</span>
                            <span className="bg-red-500 text-[8px] px-1 rounded">FREE</span>
                          </div>
                          <button className="bg-white text-gray-800 text-[10px] px-3 py-1 rounded-full font-medium">
                            Transfer
                          </button>
                        </div>
                      </div>
                      <div className="bg-gray-100 rounded-xl p-3">
                        <h4 className="text-xs font-semibold text-gray-800 mb-1">Save Pro</h4>
                        <p className="text-[9px] text-gray-500 mb-2">Save at attractive interest rate. Budget and Save.</p>
                        <div className="flex items-center justify-between">
                          <span className="text-xl">🏦</span>
                          <button className="bg-gray-800 text-white text-[10px] px-3 py-1 rounded-full font-medium">
                            Save Now
                          </button>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Race Car at bottom */}
        <div className="absolute bottom-0 left-1/2 -translate-x-1/2 w-full max-w-4xl">
          <div className="relative">
            <div className="w-full h-32 bg-gradient-to-t from-[#1a0000] to-transparent" />
            <motion.div
              initial={{ x: -100, opacity: 0 }}
              animate={{ x: 0, opacity: 1 }}
              transition={{ duration: 1.5, delay: 0.5 }}
              className="absolute bottom-0 left-1/2 -translate-x-1/2"
            >
              <div className="w-[500px] h-[120px] relative">
                {/* Stylized race car silhouette */}
                <svg viewBox="0 0 500 120" className="w-full h-full">
                  <path
                    d="M50 100 L80 80 L120 70 L180 65 L220 55 L280 50 L320 45 L380 48 L420 55 L450 65 L470 80 L480 100 Z"
                    fill="#1a0000"
                    opacity="0.8"
                  />
                  <circle cx="130" cy="100" r="25" fill="#2d1010" />
                  <circle cx="130" cy="100" r="15" fill="#4A0000" />
                  <circle cx="380" cy="100" r="25" fill="#2d1010" />
                  <circle cx="380" cy="100" r="15" fill="#4A0000" />
                  <rect x="200" y="55" width="100" height="15" rx="5" fill="#C41E3A" opacity="0.6" />
                  {/* Tail lights */}
                  <circle cx="55" cy="90" r="4" fill="#ff4444" />
                  <circle cx="65" cy="90" r="4" fill="#ff4444" />
                  <circle cx="75" cy="90" r="4" fill="#ff4444" />
                </svg>
              </div>
            </motion.div>
          </div>
        </div>
      </div>
    </section>
  );
}

// ============ FEATURES GRID ============
const features = [
  {
    icon: <Send className="w-6 h-6 text-pink-500" />,
    iconBg: 'bg-pink-50',
    title: 'Unlimited Free Transfer',
    description: 'Make seamless transfers without charges.',
  },
  {
    icon: <Receipt className="w-6 h-6 text-blue-500" />,
    iconBg: 'bg-blue-50',
    title: 'Pay Bills without Charges',
    description: 'Pay for your airtime, data plans, electricity, and TV subscription',
  },
  {
    icon: <ShoppingBag className="w-6 h-6 text-cyan-500" />,
    iconBg: 'bg-cyan-50',
    title: 'Buy Now, Pay Later',
    description: 'Buy products on the Blayzz e-shop and spread your payments.',
  },
  {
    icon: <Shield className="w-6 h-6 text-teal-500" />,
    iconBg: 'bg-teal-50',
    title: 'Premium Insurance',
    description: 'Pay your insurance premium in instalments at 0% interest rate.',
  },
  {
    icon: <Clock className="w-6 h-6 text-green-500" />,
    iconBg: 'bg-green-50',
    title: 'Budget and Save',
    description: 'Customise your savings plans, create a budget and save the unspent portion of your budget',
  },
  {
    icon: <CreditCard className="w-6 h-6 text-red-500" />,
    iconBg: 'bg-red-50',
    title: 'Pay with Blayzz card',
    description: 'Seamless, secure, swift transfer with Blayzz',
  },
  {
    icon: <Gift className="w-6 h-6 text-yellow-500" />,
    iconBg: 'bg-yellow-50',
    title: 'Get Reward',
    description: 'Earn when you transact, and when you refer family and friends.',
  },
  {
    icon: <Award className="w-6 h-6 text-purple-500" />,
    iconBg: 'bg-purple-50',
    title: 'Get Top Deals on Blayzz',
    description: 'Get great deals on blayzz from your partner merchants.',
  },
  {
    icon: <Briefcase className="w-6 h-6 text-amber-500" />,
    iconBg: 'bg-amber-50',
    title: 'Online Business',
    description: 'Get a free customised webpage on the blayzz app for your business.',
  },
  {
    icon: <Globe className="w-6 h-6 text-indigo-500" />,
    iconBg: 'bg-indigo-50',
    title: 'Integrate Blayzz to website',
    description: 'Integrate blayzz website or app for your customers to pay in instalments',
  },
];

function FeaturesGrid() {
  return (
    <section className="py-16 bg-white">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-4">
          {features.map((feature, i) => (
            <motion.div
              key={i}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ delay: i * 0.05 }}
              viewport={{ once: true }}
              className="bg-white border border-gray-200 rounded-2xl p-6 hover:shadow-lg transition-shadow duration-300 flex flex-col"
            >
              <div className={`w-12 h-12 ${feature.iconBg} rounded-xl flex items-center justify-center mb-4`}>
                {feature.icon}
              </div>
              <h3 className="font-bold text-gray-900 text-sm mb-2 leading-tight">{feature.title}</h3>
              <p className="text-xs text-gray-500 mb-6 flex-grow leading-relaxed">{feature.description}</p>
              <button className="border border-gray-300 rounded-full px-4 py-2 text-xs font-medium text-gray-700 hover:bg-gray-50 transition-colors self-start">
                Discover More
              </button>
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  );
}

// ============ CAROUSEL SECTION ============
const carouselCards = [
  {
    title: 'Transfer and pay bills',
    subtitle: 'completely free.',
    subtitleStyle: 'text-red-500 font-handwriting',
    description: 'Unlimted free transfers and pay bills at no cost',
    bg: 'bg-gradient-to-br from-rose-400 to-rose-600',
    cardTitle: 'Pay Bills Without Charges',
    cardDesc: 'Pay for your airtime, data plans, electricity, and TV subscription.',
    icon: '📄',
  },
  {
    title: 'Budget and Save',
    subtitle: 'smartly.',
    description: 'Customise your savings plans, create a budget and save the unspent portion of your budget.',
    bg: 'bg-gradient-to-br from-slate-800 to-slate-900',
    cardTitle: 'Budget and Save',
    cardDesc: 'Customise your savings plans, create a unspent portion of your budget.',
    icon: '🏦',
  },
  {
    title: 'Buy Now Pay Later',
    subtitle: 'with ease.',
    description: 'Buy products on the Blayzz e-shop and spread your payments.',
    bg: 'bg-gradient-to-br from-blue-500 to-blue-700',
    cardTitle: 'Buy Now Pay Later',
    cardDesc: 'Buy products and spread your payments conveniently.',
    icon: '🛒',
  },
];

function CarouselSection() {
  const [currentIndex, setCurrentIndex] = useState(1);

  const next = () => setCurrentIndex((prev) => (prev + 1) % carouselCards.length);
  const prev = () => setCurrentIndex((prev) => (prev - 1 + carouselCards.length) % carouselCards.length);

  return (
    <section className="py-20 bg-white">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="grid lg:grid-cols-2 gap-12 items-center">
          {/* Left Content */}
          <div className="space-y-6">
            <h2 className="text-4xl sm:text-5xl font-bold text-gray-900 leading-tight">
              {carouselCards[currentIndex].title}
              <br />
              <span className="text-red-500 italic font-serif">{carouselCards[currentIndex].subtitle}</span>
            </h2>
            <p className="text-gray-600 text-lg">{carouselCards[currentIndex].description}</p>
            
            <div className="flex gap-4 pt-4">
              <button
                onClick={prev}
                className="w-12 h-12 rounded-full border border-gray-300 flex items-center justify-center hover:bg-gray-100 transition-colors"
              >
                <ArrowLeft className="w-5 h-5 text-gray-700" />
              </button>
              <button
                onClick={next}
                className="w-12 h-12 rounded-full border border-gray-300 flex items-center justify-center hover:bg-gray-100 transition-colors"
              >
                <ArrowRight className="w-5 h-5 text-gray-700" />
              </button>
            </div>
          </div>

          {/* Right Cards */}
          <div className="relative h-[400px]">
            {carouselCards.map((card, i) => {
              const offset = i - currentIndex;
              const isActive = offset === 0;
              const isNext = offset === 1 || (currentIndex === carouselCards.length - 1 && i === 0);
              
              return (
                <motion.div
                  key={i}
                  initial={false}
                  animate={{
                    x: isActive ? 0 : isNext ? 200 : -200,
                    scale: isActive ? 1 : 0.9,
                    opacity: isActive ? 1 : 0.7,
                    zIndex: isActive ? 10 : 5,
                  }}
                  transition={{ duration: 0.4 }}
                  className={`absolute inset-0 ${card.bg} rounded-3xl p-8 text-white cursor-pointer`}
                  onClick={() => setCurrentIndex(i)}
                >
                  <div className="text-6xl mb-6">{card.icon}</div>
                  <h3 className="text-2xl font-bold mb-3">{card.cardTitle}</h3>
                  <p className="text-white/80 text-sm mb-6">{card.cardDesc}</p>
                  <button className="border border-white/50 rounded-full px-5 py-2 text-sm font-medium hover:bg-white/10 transition-colors">
                    Learn more
                  </button>
                </motion.div>
              );
            })}
          </div>
        </div>
      </div>
    </section>
  );
}

// ============ BUY NOW PAY LATER ============
const bnplTabs = [
  { label: 'Solar Panels, Inverters & Batteries', active: true },
  { label: 'Electronics, Phones & Gadgets', active: false },
  { label: 'Own your dream car', active: false },
  { label: 'Furniture & Household Equipment', active: false },
];

function BuyNowPayLater() {
  const [activeTab, setActiveTab] = useState(0);

  return (
    <section className="py-20 bg-white">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="mb-12">
          <h2 className="text-4xl font-bold text-gray-900 mb-3">Buy now pay later</h2>
          <p className="text-gray-600">With blayzz, you have no worries. Simply split your payment at your convenience.</p>
        </div>

        {/* Tabs */}
        <div className="flex flex-wrap gap-0 border-b border-gray-200 mb-12">
          {bnplTabs.map((tab, i) => (
            <button
              key={i}
              onClick={() => setActiveTab(i)}
              className={`px-6 py-3 text-sm font-medium border-b-2 transition-colors whitespace-nowrap ${
                activeTab === i
                  ? 'border-red-500 text-gray-900'
                  : 'border-transparent text-gray-500 hover:text-gray-700'
              }`}
            >
              {tab.label}
            </button>
          ))}
        </div>

        {/* Content */}
        <AnimatePresence mode="wait">
          <motion.div
            key={activeTab}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -20 }}
            className="grid lg:grid-cols-2 gap-12 items-center"
          >
            <div className="flex justify-center">
              <div className="relative w-full max-w-md">
                {/* Solar panel illustration */}
                <div className="bg-gradient-to-br from-blue-100 to-blue-200 rounded-3xl p-8 aspect-square flex items-center justify-center">
                  <div className="grid grid-cols-2 gap-4">
                    <div className="w-32 h-40 bg-blue-600 rounded-lg grid grid-cols-4 grid-rows-6 gap-1 p-2 shadow-lg">
                      {Array.from({ length: 24 }).map((_, i) => (
                        <div key={i} className="bg-blue-400/50 rounded-sm" />
                      ))}
                    </div>
                    <div className="w-32 h-40 bg-blue-700 rounded-lg grid grid-cols-4 grid-rows-6 gap-1 p-2 shadow-lg">
                      {Array.from({ length: 24 }).map((_, i) => (
                        <div key={i} className="bg-blue-500/50 rounded-sm" />
                      ))}
                    </div>
                    <div className="col-span-2 flex justify-center gap-4 mt-4">
                      <div className="w-20 h-16 bg-gray-200 rounded-lg border-2 border-gray-400 flex items-center justify-center">
                        <div className="w-2 h-2 bg-red-500 rounded-full mr-1" />
                        <div className="w-2 h-2 bg-black rounded-full" />
                      </div>
                      <div className="w-16 h-20 bg-gray-800 rounded-lg flex flex-col items-center justify-center gap-2">
                        <div className="w-8 h-4 bg-blue-400 rounded" />
                        <div className="w-1 h-8 bg-blue-400/50 rounded" />
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div className="space-y-6">
              <h3 className="text-2xl font-bold text-gray-900">
                Solar Panels, Inverters & Batteries
              </h3>
              <p className="text-gray-600 leading-relaxed">
                With blayzz, renewable energy is the key to a brighter future - let's unlock its potential!.
              </p>
              <button className="border border-gray-300 rounded-full px-6 py-3 text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors">
                Learn more
              </button>
            </div>
          </motion.div>
        </AnimatePresence>
      </div>
    </section>
  );
}

// ============ HOW IT WORKS ============
const steps = [
  { number: '01', title: 'Sign up on Blayzz', color: 'text-green-500' },
  { number: '02', title: 'Add product to cart', color: 'text-blue-500' },
  { number: '03', title: 'Pay in instalments', color: 'text-yellow-500', active: true },
  { number: '04', title: 'Easy checkout', color: 'text-red-400' },
];

function HowItWorks() {
  return (
    <section className="py-20 bg-gray-50">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        <h2 className="text-4xl font-bold text-gray-900 mb-12">
          How pay in <span className="text-red-600">instalments</span> works
        </h2>

        <div className="grid lg:grid-cols-2 gap-12 items-center">
          {/* Steps */}
          <div className="space-y-4">
            {steps.map((step, i) => (
              <motion.div
                key={i}
                initial={{ opacity: 0, x: -20 }}
                whileInView={{ opacity: 1, x: 0 }}
                transition={{ delay: i * 0.1 }}
                viewport={{ once: true }}
                className={`bg-white rounded-xl p-6 border ${
                  step.active ? 'border-gray-200 shadow-md' : 'border-gray-100'
                }`}
              >
                <div className="flex items-center gap-3 mb-3">
                  <span className={`text-2xl font-bold ${step.color}`}>{step.number}</span>
                  <h3 className={`text-xl ${step.active ? 'font-bold text-gray-900' : 'font-medium text-gray-500'}`}>
                    {step.title}
                  </h3>
                </div>
                {step.active && (
                  <p className="text-sm text-gray-600 ml-12 bg-gray-50 rounded-lg p-4">
                    When you choose to pay in instalments, you pre-qualified yourself on the app and get an instant decision.
                  </p>
                )}
              </motion.div>
            ))}
          </div>

          {/* Phone Mockup */}
          <div className="flex justify-center">
            <div className="relative">
              <div className="absolute inset-0 bg-gradient-to-br from-pink-300 to-red-500 rounded-3xl transform rotate-3 scale-105" />
              <div className="relative w-[280px] bg-gray-900 rounded-[2.5rem] p-3 shadow-2xl">
                <div className="bg-white rounded-[2rem] overflow-hidden">
                  {/* Cart Screen */}
                  <div className="bg-gray-100 px-4 py-6">
                    <div className="bg-white rounded-xl p-3 mb-3 flex items-center gap-3">
                      <div className="w-12 h-12 bg-gray-200 rounded-lg" />
                      <div className="flex-1">
                        <p className="text-xs font-medium text-gray-800">Hisense AC SPL 2HP</p>
                        <p className="text-[10px] text-gray-500">Quantity: 2</p>
                        <p className="text-xs font-bold text-gray-900">₦1,066,000.00</p>
                        <span className="text-[9px] bg-gray-200 px-2 py-0.5 rounded">Pickup location</span>
                      </div>
                      <button className="text-gray-400">🗑</button>
                    </div>
                    <div className="flex justify-between text-xs text-gray-600 py-2">
                      <span>Subtotal</span>
                      <span>₦1,066,000.00</span>
                    </div>
                    <div className="flex justify-between text-sm font-bold text-gray-900 py-2 border-t">
                      <span>Total</span>
                      <span>₦1,066,000.00</span>
                    </div>
                    <button className="w-full bg-gray-900 text-white py-3 rounded-xl mt-3 font-medium">
                      Pay
                    </button>
                  </div>

                  {/* Payment Options Modal */}
                  <div className="bg-white px-4 py-4 border-t-2 border-gray-200 rounded-b-[2rem]">
                    <p className="text-center text-xs text-gray-500 mb-1">Total Amount</p>
                    <p className="text-center text-xl font-bold text-gray-900 mb-4">₦1,066,000.00</p>
                    
                    <div className="space-y-2">
                      <button className="w-full flex items-center gap-3 bg-gray-50 rounded-xl p-3 text-left hover:bg-gray-100 transition-colors">
                        <div className="w-8 h-8 bg-red-100 rounded-full flex items-center justify-center">
                          <CreditCard className="w-4 h-4 text-red-500" />
                        </div>
                        <div>
                          <p className="text-xs font-semibold text-gray-900">Pay in Instalments</p>
                          <p className="text-[9px] text-gray-500">Enjoy flexible repayment at interest rate of 35% per annum</p>
                        </div>
                      </button>
                      <button className="w-full flex items-center gap-3 bg-gray-50 rounded-xl p-3 text-left hover:bg-gray-100 transition-colors">
                        <div className="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center">
                          <CreditCard className="w-4 h-4 text-blue-500" />
                        </div>
                        <div>
                          <p className="text-xs font-semibold text-gray-900">Pay with Card</p>
                          <p className="text-[9px] text-gray-500">Pay securely with Debit or Credit Card</p>
                        </div>
                      </button>
                      <button className="w-full flex items-center gap-3 bg-gray-50 rounded-xl p-3 text-left hover:bg-gray-100 transition-colors">
                        <div className="w-8 h-8 bg-green-100 rounded-full flex items-center justify-center">
                          <Wallet className="w-4 h-4 text-green-500" />
                        </div>
                        <div>
                          <p className="text-xs font-semibold text-gray-900">Pay with Account</p>
                          <p className="text-[9px] text-gray-500">Make seamless payments</p>
                        </div>
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

// ============ BUSINESS SECTION ============
const businessTabs = [
  'Get a free online webpage',
  'Get a free shop on blayzz',
  'Offer pay in instalments',
  'Integrate to blayzz',
];

function BusinessSection() {
  const [activeTab, setActiveTab] = useState(0);

  return (
    <section className="py-20 bg-white">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center mb-12">
          <h2 className="text-4xl font-bold text-gray-900 mb-3">Blayzz works for your business</h2>
          <p className="text-gray-600 max-w-2xl mx-auto">
            With blayzz, you can boost your business online or in-store and offer pay in instalments.
          </p>
        </div>

        {/* Tabs */}
        <div className="flex flex-wrap justify-center gap-2 bg-gray-100 rounded-full p-2 mb-12 max-w-3xl mx-auto">
          {businessTabs.map((tab, i) => (
            <button
              key={i}
              onClick={() => setActiveTab(i)}
              className={`px-5 py-2.5 rounded-full text-sm font-medium transition-all ${
                activeTab === i
                  ? 'bg-red-600 text-white shadow-md'
                  : 'text-gray-600 hover:text-gray-900'
              }`}
            >
              {tab}
            </button>
          ))}
        </div>

        {/* Content */}
        <div className="grid lg:grid-cols-2 gap-12 items-center">
          {/* Mockup Images */}
          <div className="relative">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-4">
                <div className="bg-gray-900 rounded-2xl overflow-hidden aspect-[3/4] relative">
                  <div className="absolute inset-0 bg-gradient-to-br from-purple-900 to-black" />
                  <div className="absolute bottom-4 left-4 right-4">
                    <p className="text-white font-bold text-lg">Blayzz Business</p>
                  </div>
                </div>
                <div className="bg-gradient-to-br from-orange-400 to-orange-600 rounded-2xl overflow-hidden aspect-[3/4] relative p-6">
                  <p className="text-white font-bold text-2xl mb-4">Sell more with style</p>
                  <p className="text-white/80 text-xs mb-4">CREATE CAMPAIGNS</p>
                  <button className="bg-white text-gray-900 px-4 py-2 rounded-full text-sm font-medium">
                    Explore
                  </button>
                </div>
              </div>
              <div className="space-y-4 mt-8">
                <div className="bg-gradient-to-br from-blue-500 to-purple-600 rounded-2xl overflow-hidden aspect-[3/4] relative p-6">
                  <p className="text-white font-bold text-2xl mb-2">Launch your business</p>
                  <p className="text-white/80 text-xs mb-4">Get started in 3 minutes</p>
                  <div className="bg-white/20 backdrop-blur-sm rounded-lg p-3">
                    <p className="text-white text-xs">Username</p>
                    <div className="h-6 bg-white/30 rounded mt-1" />
                    <p className="text-white text-xs mt-2">Email</p>
                    <div className="h-6 bg-white/30 rounded mt-1" />
                    <div className="flex gap-2 mt-3">
                      <div className="flex-1 h-8 bg-white rounded" />
                      <div className="flex-1 h-8 bg-orange-400 rounded" />
                    </div>
                  </div>
                </div>
                <div className="bg-gray-800 rounded-2xl overflow-hidden aspect-[3/4] relative p-4">
                  <div className="flex gap-2 mb-3">
                    {['Alcent', 'Artigras', 'Axi Invest', 'Aust', 'eiconnect', 'Commerect', 'Cultur'].map((t, i) => (
                      <span key={i} className="text-[8px] text-gray-400 whitespace-nowrap">{t}</span>
                    ))}
                  </div>
                  <div className="bg-gray-700 rounded-lg p-3">
                    <p className="text-white text-xs font-medium">Pay@Berges</p>
                    <div className="h-16 bg-gray-600 rounded mt-2" />
                  </div>
                  <p className="text-white font-bold text-lg mt-4 text-center">Buy now pay later</p>
                </div>
              </div>
            </div>
          </div>

          {/* Text Content */}
          <div className="space-y-6">
            <h3 className="text-3xl font-bold text-gray-900">Get a free online webpage</h3>
            <p className="text-gray-600 leading-relaxed">
              Do you want an online presence, blayzz can give you a customised webpage for your business so that your customers can make online purchases.
            </p>
            <button className="bg-gray-900 text-white px-8 py-3 rounded-full font-medium hover:bg-gray-800 transition-colors">
              Sign up
            </button>
          </div>
        </div>
      </div>
    </section>
  );
}

// ============ CTA SECTION ============
function CTASection() {
  return (
    <section className="py-20 bg-gray-50">
      <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="relative bg-gradient-to-br from-[#8B0000] via-[#A00000] to-[#6B0000] rounded-3xl overflow-hidden py-20 px-8 text-center">
          {/* Concentric circles */}
          <div className="absolute inset-0 flex items-center justify-center">
            {[1, 2, 3, 4, 5].map((i) => (
              <div
                key={i}
                className="absolute rounded-full border border-white/10"
                style={{
                  width: `${i * 200}px`,
                  height: `${i * 200}px`,
                }}
              />
            ))}
          </div>

          <div className="relative z-10">
            <div className="inline-flex items-center justify-center w-24 h-24 bg-white/10 backdrop-blur-sm rounded-3xl mb-8 border border-white/20">
              <div className="w-16 h-16 bg-red-600 rounded-2xl flex items-center justify-center transform rotate-12">
                <ArrowUp className="w-10 h-10 text-white" strokeWidth={3} />
              </div>
            </div>

            <h2 className="text-4xl sm:text-5xl font-bold text-white mb-4">
              Freedom to pay your way
            </h2>
            <p className="text-white/80 text-lg mb-8 max-w-2xl mx-auto">
              We're building the smartest way to spend, send, and split payment anytime, anywhere. Zero stress.
            </p>

            <div className="inline-flex items-center gap-3 bg-white rounded-full px-6 py-3 shadow-lg">
              <span className="text-gray-800 font-medium text-sm">Get the app</span>
              <Play className="w-4 h-4 text-green-600" fill="currentColor" />
              <Apple className="w-4 h-4 text-gray-800" fill="currentColor" />
              <QrCode className="w-4 h-4 text-gray-800" />
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

// ============ FOOTER ============
function Footer() {
  return (
    <footer className="bg-gray-50 border-t border-gray-200 pt-16 pb-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-12 mb-12">
          {/* Blayzz Logo & Info */}
          <div className="space-y-4">
            <div className="flex items-center gap-1">
              <div className="w-8 h-8 bg-red-600 rounded-lg flex items-center justify-center">
                <ArrowUp className="w-5 h-5 text-white" strokeWidth={3} />
              </div>
              <span className="text-gray-900 font-bold text-xl ml-1">blayzz</span>
            </div>
            <p className="text-sm text-gray-600">Tailor-made solutions that ensure a lifestyle of comfort.</p>
            <div className="flex items-center gap-2 text-sm text-gray-600">
              <span className="text-green-600">🇳🇬</span>
              <span>Licensed by the Central Bank of Nigeria</span>
            </div>
            <p className="text-xs text-gray-500">©2026 Blayzz by PremiumTrust Bank</p>
            <button className="border border-gray-300 rounded px-4 py-2 text-xs font-medium text-gray-700 hover:bg-gray-100 transition-colors">
              Review Cookie Consent
            </button>
          </div>

          {/* PremiumTrust Bank */}
          <div className="space-y-4">
            <div className="flex items-center gap-2">
              <div className="w-8 h-8 border-2 border-gray-800 rounded flex items-center justify-center">
                <ArrowUp className="w-4 h-4 text-red-600" strokeWidth={3} />
              </div>
              <div>
                <span className="text-gray-900 font-bold">Premium</span>
                <span className="text-red-600 font-bold">Trust</span>
                <br />
                <span className="text-gray-900 font-bold">Bank</span>
              </div>
            </div>
            <div className="text-sm text-gray-600 space-y-1">
              <p>1612 Adeola Hopewell Street, Victoria Island, Lagos State, Nigeria</p>
              <p>0700PREMIUM (07007736486),</p>
              <p>02013302777</p>
            </div>
          </div>

          {/* Legal */}
          <div>
            <h4 className="font-bold text-gray-900 mb-4">Legal</h4>
            <ul className="space-y-3 text-sm text-gray-600">
              <li><a href="#" className="hover:text-gray-900 transition-colors">Privacy policy</a></li>
              <li><a href="#" className="hover:text-gray-900 transition-colors">Terms & conditions</a></li>
              <li><a href="#" className="hover:text-gray-900 transition-colors">Cookie policy</a></li>
              <li><a href="#" className="hover:text-gray-900 transition-colors">ISMS policy</a></li>
            </ul>
          </div>

          {/* Help & Security */}
          <div>
            <h4 className="font-bold text-gray-900 mb-4">Help & security</h4>
            <ul className="space-y-3 text-sm text-gray-600">
              <li><a href="#" className="hover:text-gray-900 transition-colors">Contact us</a></li>
              <li><a href="#" className="hover:text-gray-900 transition-colors">Security center</a></li>
              <li><a href="#" className="hover:text-gray-900 transition-colors">FAQs</a></li>
              <li><a href="#" className="hover:text-gray-900 transition-colors">contactpremium@premiumtrustbank.com</a></li>
            </ul>
            <div className="flex gap-3 mt-6">
              <a href="#" className="w-8 h-8 bg-gray-700 rounded-full flex items-center justify-center hover:bg-gray-800 transition-colors">
                <svg className="w-4 h-4 text-white" fill="currentColor" viewBox="0 0 24 24"><path d="M12 2.163c3.204 0 3.584.012 4.85.07 3.252.148 4.771 1.691 4.919 4.919.058 1.265.069 1.645.069 4.849 0 3.205-.012 3.584-.069 4.849-.149 3.225-1.664 4.771-4.919 4.919-1.266.058-1.644.07-4.85.07-3.204 0-3.584-.012-4.849-.07-3.26-.149-4.771-1.699-4.919-4.92-.058-1.265-.07-1.644-.07-4.849 0-3.204.013-3.583.07-4.849.149-3.227 1.664-4.771 4.919-4.919 1.266-.057 1.645-.069 4.849-.069zM12 0C8.741 0 8.333.014 7.053.072 2.695.272.273 2.69.073 7.052.014 8.333 0 8.741 0 12c0 3.259.014 3.668.072 4.948.2 4.358 2.618 6.78 6.98 6.98C8.333 23.986 8.741 24 12 24c3.259 0 3.668-.014 4.948-.072 4.354-.2 6.782-2.618 6.979-6.98.059-1.28.073-1.689.073-4.948 0-3.259-.014-3.667-.072-4.947-.196-4.354-2.617-6.78-6.979-6.98C15.668.014 15.259 0 12 0zm0 5.838a6.162 6.162 0 100 12.324 6.162 6.162 0 000-12.324zM12 16a4 4 0 110-8 4 4 0 010 8zm6.406-11.845a1.44 1.44 0 100 2.881 1.44 1.44 0 000-2.881z"/></svg>
              </a>
              <a href="#" className="w-8 h-8 bg-gray-700 rounded-full flex items-center justify-center hover:bg-gray-800 transition-colors">
                <svg className="w-4 h-4 text-white" fill="currentColor" viewBox="0 0 24 24"><path d="M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.046c.477-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433a2.062 2.062 0 01-2.063-2.065 2.064 2.064 0 112.063 2.065zm1.782 13.019H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451C23.2 24 24 23.227 24 22.271V1.729C24 .774 23.2 0 22.222 0h.003z"/></svg>
              </a>
              <a href="#" className="w-8 h-8 bg-gray-700 rounded-full flex items-center justify-center hover:bg-gray-800 transition-colors">
                <svg className="w-4 h-4 text-white" fill="currentColor" viewBox="0 0 24 24"><path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z"/></svg>
              </a>
              <a href="#" className="w-8 h-8 bg-gray-700 rounded-full flex items-center justify-center hover:bg-gray-800 transition-colors">
                <svg className="w-4 h-4 text-white" fill="currentColor" viewBox="0 0 24 24"><path d="M24 12.073c0-6.627-5.373-12-12-12s-12 5.373-12 12c0 5.99 4.388 10.954 10.125 11.854v-8.385H7.078v-3.47h3.047V9.43c0-3.007 1.792-4.669 4.533-4.669 1.312 0 2.686.235 2.686.235v2.953H15.83c-1.491 0-1.956.925-1.956 1.874v2.25h3.328l-.532 3.47h-2.796v8.385C19.612 23.027 24 18.062 24 12.073z"/></svg>
              </a>
            </div>
          </div>
        </div>
      </div>
    </footer>
  );
}

// ============ SCROLL TO TOP ============
function ScrollToTop() {
  const [visible, setVisible] = useState(false);

  useEffect(() => {
    const handleScroll = () => setVisible(window.scrollY > 500);
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  return (
    <AnimatePresence>
      {visible && (
        <motion.button
          initial={{ opacity: 0, scale: 0.8 }}
          animate={{ opacity: 1, scale: 1 }}
          exit={{ opacity: 0, scale: 0.8 }}
          onClick={() => window.scrollTo({ top: 0, behavior: 'smooth' })}
          className="fixed bottom-8 right-8 w-12 h-12 bg-white border border-red-200 rounded-full shadow-lg flex items-center justify-center hover:bg-red-50 transition-colors z-50"
        >
          <ArrowUp className="w-5 h-5 text-red-600" />
        </motion.button>
      )}
    </AnimatePresence>
  );
}

// ============ MAIN APP ============
export default function App() {
  return (
    <div className="min-h-screen bg-white">
      <Navbar />
      <HeroSection />
      <FeaturesGrid />
      <CarouselSection />
      <BuyNowPayLater />
      <HowItWorks />
      <BusinessSection />
      <CTASection />
      <Footer />
      <ScrollToTop />
    </div>
  );
}
