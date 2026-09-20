import { emptySong, emptySongSection, reduxRender } from '../../../../../../test-utils.tsx'
import { SongProvider } from '../../context/SongContext.tsx'
import { SongOutlineProvider } from '../outline/context/SongOutlineContext.tsx'
import { ClickSelectProvider } from '../../../../../../context/ClickSelectContext.tsx'
import SongSections from './SongSections.tsx'
import Song, { SongSection } from '../../../../../../types/models/Song.ts'
import { fireEvent, screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { expect } from 'vitest'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { MoveSongSectionRequest } from './types/requests/SongSectionRequests.ts'
import { createRef, ReactNode } from 'react'

// Mock Main Context
vi.mock('../../../../../../context/MainContext.tsx', () => ({
  useMain: vi.fn(() => ({
    ref: createRef(),
    mainScroll: { ref: createRef() }
  }))
}))

describe('Song Sections', () => {
  const sections: SongSection[] = [
    {
      ...emptySongSection,
      id: '1',
      name: 'Chorus 1',
      rehearsals: 0,
      confidence: 0,
      progress: 0,
      songSectionType: {
        id: '',
        name: 'Chorus'
      }
    },
    {
      ...emptySongSection,
      id: '2',
      name: 'James Solo',
      rehearsals: 7,
      confidence: 50,
      progress: 163,
      songSectionType: {
        id: '',
        name: 'Solo'
      }
    },
    {
      ...emptySongSection,
      id: '3',
      name: 'James Riff',
      rehearsals: 1,
      confidence: 36,
      progress: 40,
      songSectionType: {
        id: '',
        name: 'Riff'
      }
    }
  ]

  const handlers = [
    http.get('/songs/sections/types', () => {
      return HttpResponse.json([])
    }),
    http.get('/songs/instruments', () => {
      return HttpResponse.json([])
    })
  ]

  const server = setupServer(...handlers)

  afterEach(() => {
    server.resetHandlers()
    vi.restoreAllMocks()
  })

  beforeAll(() => server.listen())

  afterAll(() => server.close())

  function render(ui: ReactNode, song: Song = emptySong) {
    return reduxRender(
      <SongProvider song={song}>
        <SongOutlineProvider>
          <ClickSelectProvider data={sections}>{ui}</ClickSelectProvider>
        </SongOutlineProvider>
      </SongProvider>
    )
  }

  it('should render', async () => {
    render(<SongSections sections={sections} />)

    sections.forEach((section) => {
      expect(screen.getByLabelText(`song-section-${section.name}`)).toBeInTheDocument()
    })
  })

  it.skip('should be able to reorder sections', async () => {
    const section = sections[0]
    const overSection = sections[2]
    const songId = 'some-id'

    let capturedRequest: MoveSongSectionRequest
    server.use(
      http.put('/songs/sections/types', async (req) => {
        capturedRequest = (await req.request.json()) as MoveSongSectionRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    render(<SongSections sections={sections} />, { ...emptySong, id: songId })

    fireEvent.mouseDown(screen.getByLabelText(`song-section-${section.name}`))
    fireEvent.dragStart(screen.getByLabelText(`song-section-${section.name}`))
    fireEvent.dragOver(screen.getByLabelText(`song-section-${overSection.name}`))
    fireEvent.drop(screen.getByLabelText(`song-section-${overSection.name}`))
    fireEvent.mouseUp(screen.getByLabelText(`song-section-${overSection.name}`))

    const renderedSections = await screen.findAllByLabelText(/song-section-(?!details)/)
    const expectedSections = [sections[1], sections[2], sections[0]]

    for (let i = 0; i < sections.length; i++) {
      expect(renderedSections[i]).toHaveAccessibleName(`song-section-${expectedSections[i].name}`)
    }

    expect(capturedRequest).toStrictEqual({
      id: section.id,
      overId: overSection.id,
      songId: songId
    })
  })

  it('should show the drawer when selecting sections, and the context menu when right-clicking after selection', async () => {
    const user = userEvent.setup()

    render(<SongSections sections={sections} />)

    // selection drawer
    await user.keyboard('{Control>}')
    await user.click(screen.getByLabelText(`song-section-${sections[0].name}`))
    await user.keyboard('{/Control}')
    expect(screen.getByLabelText('song-sections-selection-drawer')).toBeInTheDocument()

    // context menu
    // await user.pointer({
    //   keys: '[MouseRight>]',
    //   target: screen.getByLabelText(`song-section-card-${songs[0].title}`)
    // })
    // expect(await screen.findByRole('menu', { name: 'song-sections-context-menu' })).toBeInTheDocument()
  })
})
