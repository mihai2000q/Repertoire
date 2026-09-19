import {
  emptySong,
  emptySongArrangement,
  emptySongPart,
  emptySongSection,
  reduxRender
} from '../../../../../../test-utils.tsx'
import { SongProvider } from '../../context/SongContext.tsx'
import Song, { SongArrangement, SongPart, SongSection } from '../../../../../../types/models/Song.ts'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { expect } from 'vitest'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { createRef } from 'react'
import { ReactNode } from 'react'
import SongOutlineWidget from './SongOutlineWidget.tsx'
import OutlineView from './types/enums/OutlineView.ts'
import { SongOutlineProvider } from './context/SongOutlineContext.tsx'

// Mock Main Context
vi.mock('../../../../../../context/MainContext.tsx', () => ({
  useMain: vi.fn(() => ({
    ref: createRef(),
    mainScroll: { ref: createRef() }
  }))
}))

describe('Song Outline Widget', () => {
  const sections: SongSection[] = [
    {
      ...emptySongSection,
      id: '1',
      name: 'Chorus 1',
      rehearsals: 0,
      confidence: 0,
      progress: 0
    },
    {
      ...emptySongSection,
      id: '2',
      name: 'Solo',
      rehearsals: 7,
      confidence: 50,
      progress: 163
    },
    {
      ...emptySongSection,
      id: '3',
      name: 'Verse',
      rehearsals: 1,
      confidence: 36,
      progress: 40
    }
  ]
  const parts: SongPart[] = [
    {
      ...emptySongPart,
      id: '1',
      name: 'Chorus 1',
      rehearsals: 0,
      confidence: 0,
      progress: 0
    },
    {
      ...emptySongPart,
      id: '2',
      name: 'James Solo',
      rehearsals: 7,
      confidence: 50,
      progress: 163
    },
    {
      ...emptySongPart,
      id: '3',
      name: 'James Riff',
      rehearsals: 1,
      confidence: 36,
      progress: 40
    }
  ]

  const arrangements: SongArrangement[] = [{ ...emptySongArrangement, id: '1' }]

  const handlers = [
    http.get('/songs/arrangements', () => {
      return HttpResponse.json(arrangements)
    }),
    http.get('/songs/instruments', () => {
      return HttpResponse.json([])
    }),
    http.get('/songs/sections/types', () => {
      return HttpResponse.json([])
    })
  ]

  const server = setupServer(...handlers)

  afterEach(() => {
    server.resetHandlers()
    vi.restoreAllMocks()
    localStorage.clear()
  })

  beforeAll(() => server.listen())

  afterAll(() => server.close())

  function render(
    ui: ReactNode,
    song: Song = emptySong,
    initialView?: OutlineView
  ) {
    return reduxRender(
      <SongProvider song={song}>
        <SongOutlineProvider initialView={initialView}>{ui}</SongOutlineProvider>
      </SongProvider>
    )
  }

  it('should render for sections', async () => {
    server.use(
      http.get('/songs/sections', () => {
        return HttpResponse.json(sections)
      })
    )

    render(<SongOutlineWidget />)

    expect(screen.getByLabelText(/outline-loader/i)).toBeInTheDocument()
    expect(await screen.findByText(/outline/i)).toBeInTheDocument()
    expect(screen.getByLabelText('outline-toolbar')).toBeInTheDocument()
    expect(screen.getByLabelText('song-sections')).toBeInTheDocument()
  })

  it('should render for parts', async () => {
    server.use(
      http.get('/songs/parts', () => {
        return HttpResponse.json(parts)
      })
    )

    render(<SongOutlineWidget />, emptySong, OutlineView.Parts)

    expect(screen.getByLabelText(/outline-loader/i)).toBeInTheDocument()
    expect(await screen.findByText(/outline/i)).toBeInTheDocument()
    expect(screen.getByLabelText('outline-toolbar')).toBeInTheDocument()
    expect(screen.getByLabelText('song-parts')).toBeInTheDocument()
  })

  it('should show add new song section card and open add new song section, when there are no sections', async () => {
    const user = userEvent.setup()

    server.use(
      http.get('/songs/sections', () => {
        return HttpResponse.json([])
      })
    )

    render(<SongOutlineWidget />)

    expect(await screen.findByLabelText('add-new-song-section-card')).toBeInTheDocument()
    await user.click(screen.getByLabelText('add-new-song-section-card'))
    expect(screen.getByLabelText('add-new-song-section')).toBeInTheDocument()
  })

  it('should show add new song part card and open add new song part, when there are no parts', async () => {
    const user = userEvent.setup()

    server.use(
      http.get('/songs/parts', () => {
        return HttpResponse.json([])
      })
    )

    render(<SongOutlineWidget />, emptySong, OutlineView.Parts)

    expect(await screen.findByLabelText('add-new-song-part-card')).toBeInTheDocument()
    await user.click(screen.getByLabelText('add-new-song-part-card'))
    expect(screen.getByLabelText('add-new-song-part')).toBeInTheDocument()
  })
})
