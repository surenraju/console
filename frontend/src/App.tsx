import { DashboardLayout } from "@/components/dashboard/dashboard-layout"
import { DashboardPage } from "@/components/dashboard/dashboard-page"
import { ThemeProvider } from "@/components/theme-provider"

function App() {
  return (
    <ThemeProvider
      attribute="class"
      defaultTheme="system"
      enableSystem
      disableTransitionOnChange
    >
      <DashboardLayout>
        <DashboardPage />
      </DashboardLayout>
    </ThemeProvider>
  )
}

export default App
