"use client"

import { Area, AreaChart, CartesianGrid, XAxis } from "recharts"

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { IconTrendingUp } from "@tabler/icons-react"

const chartData = [
  { month: "Jan", requests: 1200, change: "+12%" },
  { month: "Feb", requests: 1400, change: "+16%" },
  { month: "Mar", requests: 1100, change: "-8%" },
  { month: "Apr", requests: 1600, change: "+45%" },
  { month: "May", requests: 1800, change: "+12%" },
  { month: "Jun", requests: 2000, change: "+11%" },
]

export function ChartArea() {
  return (
    <Card className="@container/card">
      <CardHeader>
        <CardTitle>Request Volume</CardTitle>
        <CardDescription>
          <span className="hidden @[540px]/card:block">
            Total requests for the last 6 months
          </span>
          <span className="@[540px]/card:hidden">Last 6 months</span>
        </CardDescription>
      </CardHeader>
      <CardContent className="px-2 pt-4 sm:px-6 sm:pt-6">
        <div className="aspect-auto h-[250px] w-full">
          <AreaChart data={chartData} margin={{ left: 0, right: 10 }}>
            <defs>
              <linearGradient id="fillRequests" x1="0" y1="0" x2="0" y2="1">
                <stop
                  offset="5%"
                  stopColor="var(--primary)"
                  stopOpacity={1.0}
                />
                <stop
                  offset="95%"
                  stopColor="var(--primary)"
                  stopOpacity={0.1}
                />
              </linearGradient>
            </defs>
            <CartesianGrid vertical={false} />
            <XAxis
              dataKey="month"
              tickLine={false}
              axisLine={false}
              tickMargin={8}
              minTickGap={32}
            />
            <Area
              dataKey="requests"
              type="natural"
              fill="url(#fillRequests)"
              stroke="var(--primary)"
              strokeWidth={2}
            />
          </AreaChart>
        </div>
        <div className="mt-4 flex items-center gap-2 text-sm">
          <IconTrendingUp className="h-4 w-4 text-green-500" />
          <span className="text-muted-foreground">Steady growth trend</span>
        </div>
      </CardContent>
    </Card>
  )
} 