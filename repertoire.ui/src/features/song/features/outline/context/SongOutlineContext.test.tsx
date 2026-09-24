import { act, renderHook } from '@testing-library/react'
import { ComponentProps, ReactNode } from 'react'
import { emptySong } from '../../../../../test-utils.tsx'
import Song from '../../../../../types/models/Song.ts'
import LocalStorageKeys from '../../../../../types/enums/keys/LocalStorageKeys.ts'
import OutlineView from '../types/enums/OutlineView.ts'
import { SongProvider } from '../../../context/SongContext.tsx'
import { SongOutlineProvider, useSongOutlineContext } from './SongOutlineContext.tsx'

type ProviderProps = Omit<ComponentProps<typeof SongOutlineProvider>, 'children'>

describe('Song Outline Context', () => {
  let song: Song

  beforeEach(() => {
    localStorage.clear()
    song = { ...emptySong, id: 'song-1' }
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  function setup(props: ProviderProps = {}) {
    return renderHook(() => useSongOutlineContext(), {
      wrapper: ({ children }: { children: ReactNode }) => (
        <SongProvider song={song}>
          <SongOutlineProvider {...props}>{children}</SongOutlineProvider>
        </SongProvider>
      )
    })
  }

  it('should throw when used outside the provider', () => {
    vi.spyOn(console, 'error').mockImplementation(() => {})

    expect(() => renderHook(() => useSongOutlineContext())).toThrow(
      'useSongOutlineContext must be used within SongOutlineProvider'
    )
  })

  describe('view', () => {
    it('should default to sections', () => {
      const { result } = setup()
      expect(result.current.view).toBe(OutlineView.Sections)
    })

    it('should use the initial view', () => {
      const { result } = setup({ initialView: OutlineView.Parts })
      expect(result.current.view).toBe(OutlineView.Parts)
    })

    it('should change the view', () => {
      const { result } = setup()

      act(() => {
        result.current.setView(OutlineView.Parts)
      })

      expect(result.current.view).toBe(OutlineView.Parts)
    })

    // these two depend on useLocalStorage reading storage on mount;
    // if it does that in an effect, wrap the expects in waitFor
    it('should restore the stored view', () => {
      localStorage.setItem(LocalStorageKeys.SongOutlineView, JSON.stringify(OutlineView.Parts))

      const { result } = setup()

      expect(result.current.view).toBe(OutlineView.Parts)
    })

    it('should fall back to sections when the stored view is invalid', () => {
      localStorage.setItem(LocalStorageKeys.SongOutlineView, JSON.stringify('garbage'))

      const { result } = setup({ initialView: OutlineView.Parts })

      expect(result.current.view).toBe(OutlineView.Sections)
    })
  })

  describe('details', () => {
    it('should not show details by default', () => {
      const { result } = setup()
      expect(result.current.isDetailsShowing('a')).toBe(false)
    })

    it('should toggle details for a single id', () => {
      const { result } = setup()

      act(() => {
        result.current.toggleDetails('a')
      })
      expect(result.current.isDetailsShowing('a')).toBe(true)
      expect(result.current.isDetailsShowing('b')).toBe(false)

      act(() => {
        result.current.toggleDetails('a')
      })
      expect(result.current.isDetailsShowing('a')).toBe(false)
    })

    it('should use the initial details section ids in sections view', () => {
      const { result } = setup({ initialDetailsSectionIds: new Set(['s1']) })

      expect(result.current.isDetailsShowing('s1')).toBe(true)
      expect(result.current.isDetailsShowing('s2')).toBe(false)
    })

    it('should use the initial details part ids in parts view only', () => {
      const { result } = setup({
        initialView: OutlineView.Parts,
        initialDetailsPartIds: new Set(['p1'])
      })
      expect(result.current.isDetailsShowing('p1')).toBe(true)

      act(() => {
        result.current.setView(OutlineView.Sections)
      })
      expect(result.current.isDetailsShowing('p1')).toBe(false)
    })

    it('should keep details per view when switching views', () => {
      const { result } = setup()

      act(() => {
        result.current.toggleDetails('x')
      })
      act(() => {
        result.current.setView(OutlineView.Parts)
      })
      expect(result.current.isDetailsShowing('x')).toBe(false)

      act(() => {
        result.current.toggleDetails('y')
      })
      act(() => {
        result.current.setView(OutlineView.Sections)
      })
      expect(result.current.isDetailsShowing('x')).toBe(true)
      expect(result.current.isDetailsShowing('y')).toBe(false)

      act(() => {
        result.current.setView(OutlineView.Parts)
      })
      expect(result.current.isDetailsShowing('y')).toBe(true)
    })
  })

  describe('areAllDetailsShowing', () => {
    it('should be false for no ids', () => {
      const { result } = setup()
      expect(result.current.areAllDetailsShowing([])).toBe(false)
    })

    it('should be false when none or only some are showing', () => {
      const { result } = setup()
      expect(result.current.areAllDetailsShowing(['a', 'b'])).toBe(false)

      act(() => {
        result.current.toggleDetails('a')
      })
      expect(result.current.areAllDetailsShowing(['a', 'b'])).toBe(false)
    })

    it('should be true when all are showing', () => {
      const { result } = setup()

      act(() => {
        result.current.toggleDetails('a')
        result.current.toggleDetails('b')
      })

      expect(result.current.areAllDetailsShowing(['a', 'b'])).toBe(true)
    })
  })

  describe('toggleAllDetails', () => {
    it('should show all when none are showing', () => {
      const { result } = setup()

      act(() => {
        result.current.toggleAllDetails(['a', 'b', 'c'])
      })

      expect(result.current.areAllDetailsShowing(['a', 'b', 'c'])).toBe(true)
    })

    it('should show all when only some are showing', () => {
      const { result } = setup()

      act(() => {
        result.current.toggleDetails('a')
      })
      act(() => {
        result.current.toggleAllDetails(['a', 'b', 'c'])
      })

      expect(result.current.areAllDetailsShowing(['a', 'b', 'c'])).toBe(true)
    })

    it('should hide all when all are showing', () => {
      const { result } = setup()

      act(() => {
        result.current.toggleAllDetails(['a', 'b', 'c'])
      })
      act(() => {
        result.current.toggleAllDetails(['a', 'b', 'c'])
      })

      expect(result.current.isDetailsShowing('a')).toBe(false)
      expect(result.current.isDetailsShowing('b')).toBe(false)
      expect(result.current.isDetailsShowing('c')).toBe(false)
    })

    it('should only affect the given ids', () => {
      const { result } = setup()

      act(() => {
        result.current.toggleDetails('other')
      })
      act(() => {
        result.current.toggleAllDetails(['a', 'b'])
      })
      act(() => {
        result.current.toggleAllDetails(['a', 'b'])
      })

      expect(result.current.isDetailsShowing('other')).toBe(true)
    })

    it('should do nothing for no ids', () => {
      const { result } = setup()

      act(() => {
        result.current.toggleDetails('a')
      })
      act(() => {
        result.current.toggleAllDetails([])
      })

      expect(result.current.isDetailsShowing('a')).toBe(true)
    })

    it('should only affect the current view', () => {
      const { result } = setup()

      act(() => {
        result.current.toggleAllDetails(['a', 'b'])
      })
      act(() => {
        result.current.setView(OutlineView.Parts)
      })

      expect(result.current.areAllDetailsShowing(['a', 'b'])).toBe(false)
    })
  })

  describe('song change', () => {
    it('should reset details for both views when the song changes', () => {
      const { result, rerender } = setup()

      act(() => {
        result.current.toggleDetails('s1')
      })
      act(() => {
        result.current.setView(OutlineView.Parts)
      })
      act(() => {
        result.current.toggleDetails('p1')
      })

      song = { ...song, id: 'song-2' }
      rerender()

      expect(result.current.isDetailsShowing('p1')).toBe(false)
      act(() => {
        result.current.setView(OutlineView.Sections)
      })
      expect(result.current.isDetailsShowing('s1')).toBe(false)
    })

    it('should keep details when re-rendered with the same song', () => {
      const { result, rerender } = setup()

      act(() => {
        result.current.toggleDetails('s1')
      })
      rerender()

      expect(result.current.isDetailsShowing('s1')).toBe(true)
    })

    it('should keep the view when the song changes', () => {
      const { result, rerender } = setup()

      act(() => {
        result.current.setView(OutlineView.Parts)
      })
      song = { ...song, id: 'song-2' }
      rerender()

      expect(result.current.view).toBe(OutlineView.Parts)
    })
  })
})
