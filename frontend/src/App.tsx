import { BrowserRouter, Routes, Route } from "react-router-dom"
import { DashboardLayout } from "@/components/dashboard/dashboard-layout"
import { DashboardPage } from "@/components/dashboard/dashboard-page"
import { LLMProvidersPage } from "@/components/llm-providers-page"
import { ThemeProvider } from "@/components/theme-provider"

function App() {
  return (
    <ThemeProvider
      attribute="class"
      defaultTheme="system"
      enableSystem
      disableTransitionOnChange
    >
      <BrowserRouter>
        <DashboardLayout>
          <Routes>
            <Route path="/" element={<DashboardPage />} />
            <Route path="/llm-providers" element={<LLMProvidersPage />} />
          </Routes>
        </DashboardLayout>
      </BrowserRouter>
    </ThemeProvider>
  )
}

export default App
