import { Navbar } from "@/components/layout/Navbar"
import { Footer } from "@/components/layout/Footer"
import { Hero } from "@/components/home/Hero"
import { Features } from "@/components/home/Features"
import { TechStack } from "@/components/home/TechStack"
import { CTA } from "@/components/home/CTA"

export default function Home() {
  return (
    <div className="min-h-screen bg-background font-sans antialiased">
      <Navbar />
      <main>
        <Hero />
        <Features />
        <TechStack />
        <CTA />
      </main>
      <Footer />
    </div>
  )
}
