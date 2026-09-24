import { renderHook } from '@testing-library/react'
import { ReactNode } from 'react'
import { emptyArtist, emptySong, emptySongSettings } from '../../../test-utils.tsx'
import Song from '../../../types/models/Song.ts'
import { BandMember } from '../../../types/models/Artist.ts'
import { SongProvider, useSongContext } from './SongContext.tsx'

describe('Song Context', () => {
  const bandMembers: BandMember[] = [
    { id: 'bm-1', name: 'James', roles: [{ id: '12', name: 'Voice' }] },
    { id: 'bm-2', name: 'Kirk', roles: [{ id: '123', name: 'Guitar' }] }
  ]

  let song: Song

  beforeEach(() => {
    song = { ...emptySong, id: 'song-1', settings: emptySongSettings }
  })

  const wrapper = ({ children }: { children: ReactNode }) => (
    <SongProvider song={song}>{children}</SongProvider>
  )

  it('should provide defaults when used outside a provider', () => {
    const { result } = renderHook(() => useSongContext())

    expect(result.current).toStrictEqual({
      songId: '',
      isArtistBand: false,
      artistBandMembers: undefined,
      settings: undefined,
      defaultArrangementId: undefined
    })
  })

  it('should provide song id, settings and default arrangement id', () => {
    song = { ...song, defaultArrangementId: 'arr-1' }

    const { result } = renderHook(() => useSongContext(), { wrapper })

    expect(result.current.songId).toBe('song-1')
    expect(result.current.settings).toBe(song.settings)
    expect(result.current.defaultArrangementId).toBe('arr-1')
  })

  it('should leave default arrangement id undefined when the song has none', () => {
    song = { ...song, defaultArrangementId: undefined }

    const { result } = renderHook(() => useSongContext(), { wrapper })

    expect(result.current.defaultArrangementId).toBeUndefined()
  })

  it('should expose band members when the artist is a band', () => {
    song = { ...song, artist: { ...emptyArtist, isBand: true, bandMembers } }

    const { result } = renderHook(() => useSongContext(), { wrapper })

    expect(result.current.isArtistBand).toBe(true)
    expect(result.current.artistBandMembers).toStrictEqual(bandMembers)
  })

  it('should not expose band members when the artist is not a band', () => {
    song = { ...song, artist: { ...emptyArtist, isBand: false, bandMembers } }

    const { result } = renderHook(() => useSongContext(), { wrapper })

    expect(result.current.isArtistBand).toBe(false)
    expect(result.current.artistBandMembers).toBeUndefined()
  })

  it('should not be a band and have no members when the song has no artist', () => {
    song = { ...song, artist: undefined } as unknown as Song

    const { result } = renderHook(() => useSongContext(), { wrapper })

    expect(result.current.isArtistBand).toBe(false)
    expect(result.current.artistBandMembers).toBeUndefined()
  })

  it('should update when the song changes', () => {
    const { result, rerender } = renderHook(() => useSongContext(), { wrapper })
    expect(result.current.songId).toBe('song-1')
    expect(result.current.isArtistBand).toBe(false)

    song = {
      ...song,
      id: 'song-2',
      defaultArrangementId: 'arr-2',
      artist: { ...emptyArtist, isBand: true, bandMembers }
    }
    rerender()

    expect(result.current.songId).toBe('song-2')
    expect(result.current.defaultArrangementId).toBe('arr-2')
    expect(result.current.isArtistBand).toBe(true)
    expect(result.current.artistBandMembers).toStrictEqual(bandMembers)
  })
})
