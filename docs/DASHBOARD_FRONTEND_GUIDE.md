# Dashboard API - Frontend Integration Guide

This guide provides frontend developers with everything needed to integrate the Dashboard APIs into the application.

## Base Configuration

```typescript
const BASE_URL = 'http://localhost:8080/api/v1';

// Default headers for all requests
const headers = {
  'Content-Type': 'application/json',
  'Authorization': `Bearer ${accessToken}`
};
```

---

## API Endpoints Summary

| Endpoint | Method | Description | Access |
|----------|--------|-------------|--------|
| `/dashboard/statistics` | GET | KPI statistics (role-based) | All authenticated |
| `/dashboard/announcements` | GET | Company announcements | All authenticated |
| `/dashboard/announcements/:id/read` | POST | Mark announcement read | All authenticated |
| `/dashboard/events` | GET | Birthdays, anniversaries | HR/Admin |
| `/dashboard/my-activity` | GET | Employee's activities | All authenticated |
| `/dashboard/quick-actions` | GET | Quick action links | All authenticated |
| `/dashboard/pending-approvals` | GET | Pending approvals | HR/Admin |

---

## 1. Dashboard Statistics

Returns different data based on user role (admin vs employee).

### Request

```typescript
// GET /dashboard/statistics
const getStatistics = async (view?: 'admin' | 'employee', date?: string) => {
  const params = new URLSearchParams();
  if (view) params.append('view', view);
  if (date) params.append('date', date); // Format: YYYY-MM-DD
  
  const response = await fetch(`${BASE_URL}/dashboard/statistics?${params}`, {
    headers
  });
  return response.json();
};
```

### Admin Response

```typescript
interface AdminStatisticsResponse {
  success: boolean;
  message: string;
  data: {
    view: 'admin';
    date: string;
    stats: {
      totalHeadcount: {
        value: number;
        change: number;
        changePercent: number;
        trend: 'up' | 'down' | 'stable';
      };
      presentToday: {
        value: number;
        total: number;
        percentage: number;
        absent: number;
        onLeave: number;
      };
      pendingApprovals: {
        value: number;
        breakdown: {
          leaveRequests: number;
          timesheetCorrections: number;
          expenseClaims: number;
        };
      };
      payrollStatus: {
        status: 'Active' | 'Processing' | 'Completed' | 'Pending';
        currentPeriod: string;
        processingDate: string;
        totalAmount: number;
        currency: string;
      };
    };
    summary: {
      newHiresThisMonth: number;
      terminationsThisMonth: number;
      openPositions: number;
      avgAttendanceRate: number;
    };
  };
}
```

### Employee Response

```typescript
interface EmployeeStatisticsResponse {
  success: boolean;
  message: string;
  data: {
    view: 'employee';
    date: string;
    employeeId: number;
    employeeName: string;
    stats: {
      leaveBalance: {
        value: number;
        unit: string;
        breakdown: {
          annual: number;
          sick: number;
          personal: number;
        };
        usedThisYear: number;
      };
      hoursThisWeek: {
        value: number;
        target: number;
        overtime: number;
        status: 'on_track' | 'behind' | 'ahead';
      };
      pendingRequests: {
        value: number;
        breakdown: {
          leaveRequests: number;
          letterRequests: number;
        };
      };
      nextPayday: {
        date: string;
        daysRemaining: number;
        estimatedNetPay: number;
        currency: string;
      };
    };
    attendance: {
      todayStatus: 'Present' | 'Absent' | 'On Leave' | 'Not Clocked In';
      checkInTime: string | null;
      checkOutTime: string | null;
      workedHoursToday: number;
    };
  };
}
```

### React Hook Example

```typescript
import { useState, useEffect } from 'react';

export const useDashboardStats = (view?: 'admin' | 'employee') => {
  const [stats, setStats] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchStats = async () => {
      try {
        setLoading(true);
        const params = view ? `?view=${view}` : '';
        const response = await fetch(`${BASE_URL}/dashboard/statistics${params}`, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('accessToken')}`
          }
        });
        const data = await response.json();
        if (data.success) {
          setStats(data.data);
        } else {
          setError(data.message);
        }
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchStats();
  }, [view]);

  return { stats, loading, error };
};
```

---

## 2. Announcements

### Fetch Announcements

```typescript
interface AnnouncementsParams {
  page?: number;
  page_size?: number;
  type?: 'info' | 'warning' | 'success' | 'urgent';
  active_only?: boolean;
}

const getAnnouncements = async (params: AnnouncementsParams = {}) => {
  const searchParams = new URLSearchParams();
  if (params.page) searchParams.append('page', params.page.toString());
  if (params.page_size) searchParams.append('page_size', params.page_size.toString());
  if (params.type) searchParams.append('type', params.type);
  if (params.active_only !== undefined) searchParams.append('active_only', params.active_only.toString());
  
  const response = await fetch(`${BASE_URL}/dashboard/announcements?${searchParams}`, {
    headers
  });
  return response.json();
};
```

### Response Type

```typescript
interface Announcement {
  id: number;
  title: string;
  description: string;
  type: 'info' | 'warning' | 'success' | 'urgent';
  priority: 'low' | 'normal' | 'high' | 'critical';
  createdAt: string;
  createdBy: {
    id: number;
    name: string;
  };
  expiresAt: string | null;
  isRead: boolean;
  attachments: Array<{
    id: number;
    name: string;
    url: string;
    size: string;
  }>;
  relativeTime: string;
}

interface AnnouncementsResponse {
  success: boolean;
  message: string;
  data: {
    announcements: Announcement[];
    meta: {
      page: number;
      pageSize: number;
      total: number;
      totalPages: number;
      unreadCount: number;
    };
  };
}
```

### Mark as Read

```typescript
const markAnnouncementAsRead = async (announcementId: number) => {
  const response = await fetch(`${BASE_URL}/dashboard/announcements/${announcementId}/read`, {
    method: 'POST',
    headers
  });
  return response.json();
};
```

### React Component Example

```tsx
const AnnouncementCard: React.FC<{ announcement: Announcement }> = ({ announcement }) => {
  const [isRead, setIsRead] = useState(announcement.isRead);

  const handleClick = async () => {
    if (!isRead) {
      await markAnnouncementAsRead(announcement.id);
      setIsRead(true);
    }
  };

  const typeColors = {
    info: 'bg-blue-100 border-blue-500',
    warning: 'bg-yellow-100 border-yellow-500',
    success: 'bg-green-100 border-green-500',
    urgent: 'bg-red-100 border-red-500'
  };

  return (
    <div 
      className={`p-4 border-l-4 rounded ${typeColors[announcement.type]} ${!isRead ? 'font-semibold' : ''}`}
      onClick={handleClick}
    >
      <div className="flex justify-between items-start">
        <h3 className="text-lg">{announcement.title}</h3>
        <span className="text-sm text-gray-500">{announcement.relativeTime}</span>
      </div>
      <p className="mt-2 text-gray-700">{announcement.description}</p>
      {announcement.attachments.length > 0 && (
        <div className="mt-2">
          {announcement.attachments.map(att => (
            <a key={att.id} href={att.url} className="text-blue-600 mr-2">
              {att.name} ({att.size})
            </a>
          ))}
        </div>
      )}
    </div>
  );
};
```

---

## 3. Upcoming Events (Admin/HR)

### Request

```typescript
interface EventsParams {
  days_ahead?: number;  // Default: 30, max: 90
  event_type?: 'birthday' | 'anniversary' | 'probation_end' | 'all';
  department_id?: number;
  page?: number;
  page_size?: number;
}

const getUpcomingEvents = async (params: EventsParams = {}) => {
  const searchParams = new URLSearchParams();
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined) searchParams.append(key, value.toString());
  });
  
  const response = await fetch(`${BASE_URL}/dashboard/events?${searchParams}`, {
    headers
  });
  return response.json();
};
```

### Response Type

```typescript
interface UpcomingEvent {
  id: number;
  employeeId: number;
  employeeName: string;
  employeePhoto: string | null;
  department: string;
  eventType: 'birthday' | 'anniversary' | 'probation_end';
  eventDate: string;
  displayDate: string;
  daysUntil: number;
  age?: number;
  yearsOfService?: number;
  probationMonths?: number;
  initials: string;
}

interface EventsResponse {
  success: boolean;
  message: string;
  data: {
    events: UpcomingEvent[];
    summary: {
      birthdaysThisMonth: number;
      anniversariesThisMonth: number;
      probationEndsThisMonth: number;
    };
    meta: {
      page: number;
      pageSize: number;
      total: number;
      totalPages: number;
    };
  };
}
```

### React Component Example

```tsx
const EventsList: React.FC = () => {
  const [events, setEvents] = useState<UpcomingEvent[]>([]);
  const [summary, setSummary] = useState(null);

  useEffect(() => {
    const fetchEvents = async () => {
      const response = await getUpcomingEvents({ days_ahead: 30 });
      if (response.success) {
        setEvents(response.data.events);
        setSummary(response.data.summary);
      }
    };
    fetchEvents();
  }, []);

  const getEventIcon = (type: string) => {
    switch (type) {
      case 'birthday': return '🎂';
      case 'anniversary': return '🎉';
      case 'probation_end': return '✅';
      default: return '📅';
    }
  };

  return (
    <div className="space-y-4">
      {summary && (
        <div className="grid grid-cols-3 gap-4 mb-4">
          <div className="bg-pink-50 p-3 rounded">
            <div className="text-2xl font-bold">{summary.birthdaysThisMonth}</div>
            <div className="text-sm">Birthdays this month</div>
          </div>
          <div className="bg-blue-50 p-3 rounded">
            <div className="text-2xl font-bold">{summary.anniversariesThisMonth}</div>
            <div className="text-sm">Anniversaries</div>
          </div>
          <div className="bg-green-50 p-3 rounded">
            <div className="text-2xl font-bold">{summary.probationEndsThisMonth}</div>
            <div className="text-sm">Probation ends</div>
          </div>
        </div>
      )}
      
      {events.map(event => (
        <div key={event.id} className="flex items-center p-3 bg-white rounded shadow">
          <div className="w-12 h-12 rounded-full bg-gray-200 flex items-center justify-center mr-4">
            {event.employeePhoto ? (
              <img src={event.employeePhoto} alt="" className="w-12 h-12 rounded-full" />
            ) : (
              <span className="text-lg font-semibold">{event.initials}</span>
            )}
          </div>
          <div className="flex-1">
            <div className="font-semibold">{event.employeeName}</div>
            <div className="text-sm text-gray-500">{event.department}</div>
          </div>
          <div className="text-right">
            <div className="text-lg">{getEventIcon(event.eventType)}</div>
            <div className="text-sm">{event.displayDate}</div>
            <div className="text-xs text-gray-500">
              {event.daysUntil === 0 ? 'Today!' : `In ${event.daysUntil} days`}
            </div>
          </div>
        </div>
      ))}
    </div>
  );
};
```

---

## 4. My Activity (Employee)

### Request

```typescript
interface ActivityParams {
  page?: number;
  page_size?: number;
  activity_type?: 'leave' | 'timesheet' | 'payslip' | 'request' | 'all';
  days?: number;  // Default: 30
}

const getMyActivity = async (params: ActivityParams = {}) => {
  const searchParams = new URLSearchParams();
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined) searchParams.append(key, value.toString());
  });
  
  const response = await fetch(`${BASE_URL}/dashboard/my-activity?${searchParams}`, {
    headers
  });
  return response.json();
};
```

### Response Type

```typescript
interface RecentActivity {
  id: number;
  type: 'leave_request' | 'timesheet_update' | 'payslip_available' | 'letter_request' | 'approval';
  title: string;
  description: string;
  status: 'pending' | 'approved' | 'rejected' | 'completed' | 'new';
  statusLabel: string;
  icon: string;
  color: string;
  createdAt: string;
  relativeTime: string;
  metadata: Record<string, any>;
  actionUrl: string;
}

interface MyActivityResponse {
  success: boolean;
  message: string;
  data: {
    activities: RecentActivity[];
    summary: {
      pendingRequests: number;
      recentApprovals: number;
      newNotifications: number;
    };
    meta: {
      page: number;
      pageSize: number;
      total: number;
      totalPages: number;
    };
  };
}
```

---

## 5. Quick Actions

Quick actions are personalized based on user role.

### Request

```typescript
const getQuickActions = async () => {
  const response = await fetch(`${BASE_URL}/dashboard/quick-actions`, {
    headers
  });
  return response.json();
};
```

### Response Type

```typescript
interface QuickActionBadge {
  count?: number;
  text?: string;
  type: 'info' | 'warning' | 'success' | 'error';
}

interface QuickAction {
  id: string;
  title: string;
  description: string;
  icon: string;
  color: string;
  href: string;
  badge: QuickActionBadge | null;
  enabled: boolean;
}

interface QuickActionsResponse {
  success: boolean;
  message: string;
  data: {
    view: 'admin' | 'employee';
    actions: QuickAction[];
  };
}
```

### React Component Example

```tsx
const QuickActions: React.FC = () => {
  const [actions, setActions] = useState<QuickAction[]>([]);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchActions = async () => {
      const response = await getQuickActions();
      if (response.success) {
        setActions(response.data.actions);
      }
    };
    fetchActions();
  }, []);

  const iconMap: Record<string, React.ReactNode> = {
    'user-plus': <UserPlusIcon />,
    'chart': <ChartIcon />,
    'check-circle': <CheckCircleIcon />,
    'currency': <CurrencyIcon />,
    'calendar': <CalendarIcon />,
    'document': <DocumentIcon />,
    'mail': <MailIcon />,
    'users': <UsersIcon />
  };

  return (
    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
      {actions.map(action => (
        <button
          key={action.id}
          onClick={() => navigate(action.href)}
          disabled={!action.enabled}
          className={`p-4 rounded-lg text-left transition hover:shadow-lg ${
            action.enabled ? 'bg-white cursor-pointer' : 'bg-gray-100 cursor-not-allowed'
          }`}
        >
          <div className={`w-10 h-10 rounded-full flex items-center justify-center mb-2 bg-${action.color}-100`}>
            {iconMap[action.icon]}
          </div>
          <div className="font-semibold flex items-center gap-2">
            {action.title}
            {action.badge && (
              <span className={`px-2 py-0.5 text-xs rounded-full bg-${action.badge.type === 'warning' ? 'yellow' : 'blue'}-500 text-white`}>
                {action.badge.count ?? action.badge.text}
              </span>
            )}
          </div>
          <div className="text-sm text-gray-500">{action.description}</div>
        </button>
      ))}
    </div>
  );
};
```

---

## 6. Pending Approvals (Admin/HR)

### Request

```typescript
interface ApprovalsParams {
  type?: 'leave' | 'timesheet' | 'expense' | 'all';
  page?: number;
  page_size?: number;
}

const getPendingApprovals = async (params: ApprovalsParams = {}) => {
  const searchParams = new URLSearchParams();
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined) searchParams.append(key, value.toString());
  });
  
  const response = await fetch(`${BASE_URL}/dashboard/pending-approvals?${searchParams}`, {
    headers
  });
  return response.json();
};
```

### Response Type

```typescript
interface PendingApproval {
  id: number;
  type: 'leave' | 'timesheet' | 'expense';
  employeeId: number;
  employeeName: string;
  employeePhoto: string | null;
  department: string;
  title: string;
  description: string;
  submittedAt: string;
  relativeTime: string;
  priority: 'low' | 'normal' | 'high' | 'urgent';
  actionUrl: string;
  // Type-specific fields
  dates?: string;
  days?: number;
  amount?: number;
  currency?: string;
  originalTime?: string | null;
  correctedTime?: string;
}

interface PendingApprovalsResponse {
  success: boolean;
  message: string;
  data: {
    summary: {
      total: number;
      leaveRequests: number;
      timesheetCorrections: number;
      expenseClaims: number;
    };
    approvals: PendingApproval[];
    meta: {
      page: number;
      pageSize: number;
      total: number;
      totalPages: number;
    };
  };
}
```

---

## Complete Dashboard Page Example

```tsx
import React, { useEffect, useState } from 'react';
import { useDashboardStats } from './hooks/useDashboardStats';

const Dashboard: React.FC = () => {
  const { stats, loading, error } = useDashboardStats();
  const [announcements, setAnnouncements] = useState([]);
  const [quickActions, setQuickActions] = useState([]);

  useEffect(() => {
    Promise.all([
      fetch(`${BASE_URL}/dashboard/announcements`, { headers }).then(r => r.json()),
      fetch(`${BASE_URL}/dashboard/quick-actions`, { headers }).then(r => r.json())
    ]).then(([announcementsRes, actionsRes]) => {
      if (announcementsRes.success) setAnnouncements(announcementsRes.data.announcements);
      if (actionsRes.success) setQuickActions(actionsRes.data.actions);
    });
  }, []);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;

  const isAdmin = stats?.view === 'admin';

  return (
    <div className="p-6 space-y-6">
      <h1 className="text-2xl font-bold">
        {isAdmin ? 'Admin Dashboard' : `Welcome, ${stats?.employeeName}`}
      </h1>

      {/* Stats Cards */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        {isAdmin ? (
          <>
            <StatCard 
              title="Total Headcount" 
              value={stats.stats.totalHeadcount.value}
              trend={stats.stats.totalHeadcount.trend}
              change={stats.stats.totalHeadcount.changePercent}
            />
            <StatCard 
              title="Present Today" 
              value={`${stats.stats.presentToday.value}/${stats.stats.presentToday.total}`}
              subtitle={`${stats.stats.presentToday.percentage.toFixed(1)}%`}
            />
            <StatCard 
              title="Pending Approvals" 
              value={stats.stats.pendingApprovals.value}
              alert={stats.stats.pendingApprovals.value > 0}
            />
            <StatCard 
              title="Payroll Status" 
              value={stats.stats.payrollStatus.status}
              subtitle={stats.stats.payrollStatus.currentPeriod}
            />
          </>
        ) : (
          <>
            <StatCard 
              title="Leave Balance" 
              value={`${stats.stats.leaveBalance.value} days`}
            />
            <StatCard 
              title="Hours This Week" 
              value={`${stats.stats.hoursThisWeek.value}/${stats.stats.hoursThisWeek.target}`}
            />
            <StatCard 
              title="Pending Requests" 
              value={stats.stats.pendingRequests.value}
            />
            <StatCard 
              title="Next Payday" 
              value={stats.stats.nextPayday.date}
              subtitle={`${stats.stats.nextPayday.daysRemaining} days`}
            />
          </>
        )}
      </div>

      {/* Quick Actions */}
      <section>
        <h2 className="text-lg font-semibold mb-4">Quick Actions</h2>
        <QuickActions actions={quickActions} />
      </section>

      {/* Announcements */}
      <section>
        <h2 className="text-lg font-semibold mb-4">
          Announcements 
          {announcements.meta?.unreadCount > 0 && (
            <span className="ml-2 px-2 py-1 bg-red-500 text-white text-xs rounded-full">
              {announcements.meta.unreadCount} new
            </span>
          )}
        </h2>
        <div className="space-y-3">
          {announcements.slice(0, 3).map(ann => (
            <AnnouncementCard key={ann.id} announcement={ann} />
          ))}
        </div>
      </section>
    </div>
  );
};

// Reusable StatCard component
const StatCard: React.FC<{
  title: string;
  value: string | number;
  subtitle?: string;
  trend?: 'up' | 'down' | 'stable';
  change?: number;
  alert?: boolean;
}> = ({ title, value, subtitle, trend, change, alert }) => (
  <div className={`bg-white p-4 rounded-lg shadow ${alert ? 'border-l-4 border-orange-500' : ''}`}>
    <div className="text-sm text-gray-500">{title}</div>
    <div className="text-2xl font-bold flex items-center gap-2">
      {value}
      {trend && (
        <span className={trend === 'up' ? 'text-green-500' : trend === 'down' ? 'text-red-500' : 'text-gray-500'}>
          {trend === 'up' ? '↑' : trend === 'down' ? '↓' : '→'}
          {change && ` ${change.toFixed(1)}%`}
        </span>
      )}
    </div>
    {subtitle && <div className="text-sm text-gray-400">{subtitle}</div>}
  </div>
);

export default Dashboard;
```

---

## Error Handling

All API responses follow this pattern:

```typescript
// Success response
{
  success: true,
  message: "Data retrieved successfully",
  data: { ... }
}

// Error response
{
  success: false,
  message: "Error description",
  error: {
    code: "ERROR_CODE",
    details: "Detailed error message"
  }
}
```

### Common Error Codes

| HTTP Status | Code | Description |
|-------------|------|-------------|
| 400 | BAD_REQUEST | Invalid request parameters |
| 401 | UNAUTHORIZED | Missing or invalid token |
| 403 | FORBIDDEN | Insufficient permissions |
| 404 | NOT_FOUND | Resource not found |
| 500 | SERVER_ERROR | Internal server error |

### Error Handling Example

```typescript
const fetchWithErrorHandling = async (url: string) => {
  try {
    const response = await fetch(url, { headers });
    const data = await response.json();
    
    if (!data.success) {
      // Handle API error
      console.error('API Error:', data.message);
      throw new Error(data.message);
    }
    
    return data;
  } catch (error) {
    // Handle network error
    if (error.name === 'TypeError') {
      throw new Error('Network error - please check your connection');
    }
    throw error;
  }
};
```

---

## Notes for Frontend Development

1. **Authentication**: Always include the Bearer token in the Authorization header
2. **Role Detection**: The `/dashboard/statistics` endpoint returns `view: 'admin'` or `view: 'employee'` - use this to conditionally render UI
3. **Pagination**: Most list endpoints support `page` and `page_size` query parameters
4. **Real-time Updates**: Consider polling `/dashboard/statistics` every 30-60 seconds for live data
5. **Caching**: Announcements and quick actions can be cached for 5 minutes
6. **Date Format**: All dates from API are ISO 8601 format - use `new Date(dateString)` to parse

---

## Changelog

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 2026-01-23 | Initial frontend guide |
