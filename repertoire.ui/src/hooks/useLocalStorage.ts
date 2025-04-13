import { useState } from 'react'
import { useDidUpdate } from '@mantine/hooks'

export function useLocalStorage<T>({
  key,
  defaultValue
}: {
  key: string
  defaultValue?: T
}): [T, (value: T) => void] {
  const [item, setItem] = useState((JSON.parse(localStorage.getItem(key)) as T) ?? defaultValue)

  useDidUpdate(() => {
    if (item !== undefined) localStorage.setItem(key, JSON.stringify(item))
    else localStorage.removeItem(key)
  }, [item])

  return [item, setItem]
}
