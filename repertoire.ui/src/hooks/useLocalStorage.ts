import { useState } from 'react'
import { useDidUpdate } from '@mantine/hooks'

interface useLocalStorageOptions<T> {
  key: string
  defaultValue?: T
  serialize?: (value: T) => string
  deserialize?: (item: string) => T
}

export default function useLocalStorage<T>({
  key,
  defaultValue,
  serialize,
  deserialize
}: useLocalStorageOptions<T>): [T, (value: T) => void] {
  const [item, setItem] = useState<T>(() => {
    try {
      const item = localStorage.getItem(key)
      if (item === null || item === undefined) return defaultValue
      return deserialize ? deserialize(item) : JSON.parse(item)
    } catch {
      // Malformed JSON, a throwing deserialize, or storage being unavailable.
      // Fall back to the default rather than crashing the render.
      return defaultValue
    }
  })

  useDidUpdate(() => {
    try {
      if (item !== undefined)
        localStorage.setItem(key, serialize ? serialize(item) : JSON.stringify(item))
      else localStorage.removeItem(key)
    } catch {
      // Quota exceeded, storage disabled, or a throwing serialize.
      // In-memory state still works; the value just won't persist.
    }
  }, [item])

  return [item, setItem]
}
