'use client';

import { useEffect, useState } from 'react';
import { fetchUser } from '@/lib/api';
import { useRouter } from 'next/navigation';

interface User {
  id: string;
  display_name: string;
  email: string;
}

export function useUser() {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const router = useRouter();

  useEffect(() => {
    fetchUser()
      .then(setUser)
      .catch(() => {
        router.push('/');
      })
      .finally(() => setLoading(false));
  }, [router]);

  return { user, loading };
}
