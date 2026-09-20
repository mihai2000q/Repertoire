import {
  emptySong,
  emptySongArrangement,
  emptySongPart,
  emptySongSection,
  emptySongSettings,
  reduxRender,
  withToastify
} from '../../../../../test-utils.tsx'
import { SongArrangement, SongPart, SongSection } from '../../../../../types/models/Song.ts'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { expect } from 'vitest'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { AddPerfectSongRehearsalRequest } from '../../../../../types/requests/SongRequests.ts'
import { createRef, ReactNode } from 'react'
import SongOutlineToolbar from './SongOutlineToolbar.tsx'
import OutlineView from '../types/enums/OutlineView.ts'
import { SongProvider } from '../../../context/SongContext.tsx'
import { SongOutlineProvider } from '../context/SongOutlineContext.tsx'

// Mock Main Context
vi.mock('../../../../../../../context/MainContext.tsx', () => ({
  useMain: vi.fn(() => ({
    ref: createRef(),
    mainScroll: { ref: createRef() }
  }))
}))

describe('Song Outline Toolbar', () => {
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
  const sections: SongSection[] = [{ ...emptySongSection, id: 's1', parts }]

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
  })

  beforeAll(() => server.listen())

  afterAll(() => server.close())

  function render(ui: ReactNode, song = emptySong, initialView = OutlineView.Sections) {
    return reduxRender(
      <SongProvider song={song}>
        <SongOutlineProvider initialView={initialView}>{ui}</SongOutlineProvider>
      </SongProvider>
    )
  }

  it('should render for sections', async () => {
    render(<SongOutlineToolbar toggleAdd={vi.fn()} sections={sections} />, {
      ...emptySong,
      defaultArrangementId: '1',
      settings: emptySongSettings
    })

    expect(screen.getByRole('button', { name: 'add-new-section' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'show-details' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'show-details' })).not.toBeDisabled()
    expect(screen.getByRole('button', { name: 'manage-song-arrangements' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'add-custom-rehearsal' })).toBeInTheDocument()
    expect(await screen.findByRole('button', { name: 'add-custom-rehearsal' })).not.toBeDisabled()
    expect(screen.getByRole('button', { name: 'add-perfect-rehearsal' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'add-perfect-rehearsal' })).not.toBeDisabled()
    expect(screen.getByRole('button', { name: 'settings' })).toBeInTheDocument()
    expect(screen.getByRole('radio', { name: 'sections-view' })).toBeInTheDocument()
    expect(screen.getByRole('radio', { name: 'sections-view' })).toBeChecked()
    expect(screen.getByRole('radio', { name: 'parts-view' })).toBeInTheDocument()
  })

  it('should render for parts', async () => {
    const user = userEvent.setup()

    render(
      <SongOutlineToolbar toggleAdd={vi.fn()} parts={parts} />,
      { ...emptySong, defaultArrangementId: '1', settings: emptySongSettings },
      OutlineView.Parts
    )

    await user.click(screen.getByRole('radio', { name: 'parts-view' }))

    expect(screen.getByRole('button', { name: 'add-new-part' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'show-details' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'show-details' })).not.toBeDisabled()
    expect(screen.getByRole('button', { name: 'manage-song-arrangements' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'add-custom-rehearsal' })).toBeInTheDocument()
    expect(await screen.findByRole('button', { name: 'add-custom-rehearsal' })).not.toBeDisabled()
    expect(screen.getByRole('button', { name: 'add-perfect-rehearsal' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'add-perfect-rehearsal' })).not.toBeDisabled()
    expect(screen.getByRole('button', { name: 'settings' })).toBeInTheDocument()
    expect(screen.getByRole('radio', { name: 'sections-view' })).toBeInTheDocument()
    expect(screen.getByRole('radio', { name: 'parts-view' })).toBeInTheDocument()
    expect(screen.getByRole('radio', { name: 'parts-view' })).toBeChecked()
  })

  it('should disable a few options when there are no sections', () => {
    render(
      <SongOutlineToolbar
        toggleAdd={vi.fn()}
        sections={[]}
      />,
      {
        ...emptySong,
        defaultArrangementId: '1',
        settings: emptySongSettings
      }
    )

    expect(screen.getByRole('button', { name: 'show-details' })).toBeDisabled()
    expect(screen.getByRole('button', { name: 'add-custom-rehearsal' })).toBeDisabled()
    expect(screen.getByRole('button', { name: 'add-perfect-rehearsal' })).toBeDisabled()
  })

  it('should disable a few options when there are no section parts', () => {
    render(
      <SongOutlineToolbar
        toggleAdd={vi.fn()}
        sections={sections.map((s) => ({ ...s, parts: [] }))}
      />,
      {
        ...emptySong,
        defaultArrangementId: '1',
        settings: emptySongSettings
      }
    )

    expect(screen.getByRole('button', { name: 'show-details' })).not.toBeDisabled()
    expect(screen.getByRole('button', { name: 'add-custom-rehearsal' })).toBeDisabled()
    expect(screen.getByRole('button', { name: 'add-perfect-rehearsal' })).toBeDisabled()
  })

  it('should disable a few options when there are no parts', () => {
    render(
      <SongOutlineToolbar toggleAdd={vi.fn()} parts={[]} />,
      { ...emptySong, defaultArrangementId: '1', settings: emptySongSettings },
      OutlineView.Parts
    )

    expect(screen.getByRole('button', { name: 'show-details' })).toBeDisabled()
    expect(screen.getByRole('button', { name: 'add-custom-rehearsal' })).toBeDisabled()
    expect(screen.getByRole('button', { name: 'add-perfect-rehearsal' })).toBeDisabled()
  })

  it('should disable perfect rehearsal when a default arrangement is not set', () => {
    render(<SongOutlineToolbar toggleAdd={vi.fn()} />, {
      ...emptySong,
      defaultArrangementId: undefined,
      settings: emptySongSettings
    })

    expect(screen.getByRole('button', { name: 'add-perfect-rehearsal' })).toBeDisabled()
  })

  it('should call toggle add when clicking on add new section button', async () => {
    const user = userEvent.setup()

    const toggleAdd = vi.fn()

    render(<SongOutlineToolbar toggleAdd={toggleAdd} />, {
      ...emptySong,
      settings: emptySongSettings
    })

    await user.click(screen.getByRole('button', { name: 'add-new-section' }))
    expect(toggleAdd).toHaveBeenCalledOnce()
  })

  it('should call toggle add when clicking on add new part button', async () => {
    const user = userEvent.setup()

    const toggleAdd = vi.fn()

    render(
      <SongOutlineToolbar toggleAdd={toggleAdd} />,
      { ...emptySong, settings: emptySongSettings },
      OutlineView.Parts
    )

    await user.click(screen.getByRole('button', { name: 'add-new-part' }))
    expect(toggleAdd).toHaveBeenCalledOnce()
  })

  it('should show details when clicking on show details', async () => {
    const user = userEvent.setup()

    render(<SongOutlineToolbar toggleAdd={vi.fn()} sections={sections} />, {
      ...emptySong,
      settings: emptySongSettings
    })

    await user.click(screen.getByRole('button', { name: 'show-details' }))
    expect(screen.queryByRole('button', { name: 'show-details' })).not.toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'hide-details' })).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: 'hide-details' }))
    expect(screen.getByRole('button', { name: 'show-details' })).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'hide-details' })).not.toBeInTheDocument()
  })

  it('should open song arrangements modal when clicking on song arrangements button', async () => {
    const user = userEvent.setup()

    render(<SongOutlineToolbar toggleAdd={vi.fn()} />, {
      ...emptySong,
      settings: emptySongSettings
    })

    await user.click(screen.getByRole('button', { name: /song-arrangements/i }))
    expect(await screen.findByRole('dialog', { name: /song arrangements/i })).toBeInTheDocument()
  })

  it('should open add perfect rehearsal popover when on clicking add perfect rehearsal button and send request', async () => {
    const user = userEvent.setup()

    let capturedRequest: AddPerfectSongRehearsalRequest
    server.use(
      http.post('/songs/perfect-rehearsal', async (req) => {
        capturedRequest = (await req.request.json()) as AddPerfectSongRehearsalRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const songId = 'some-id'

    render(withToastify(<SongOutlineToolbar toggleAdd={vi.fn()} sections={sections} />), {
      ...emptySong,
      id: songId,
      settings: emptySongSettings,
      defaultArrangementId: '1'
    })

    await user.click(screen.getByRole('button', { name: 'add-perfect-rehearsal' }))

    expect(await screen.findByRole('dialog')).toBeInTheDocument()
    expect(screen.getByText(/increase parts' rehearsals .* occurrences/i)).toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: 'confirm' }))

    expect(await screen.findByText(/perfect rehearsal added/i)).toBeInTheDocument()
    expect(capturedRequest).toStrictEqual({ id: songId })
  })

  it('should hide details when song changes', async () => {
    const user = userEvent.setup()

    const [{ rerender }] = render(<SongOutlineToolbar toggleAdd={vi.fn()} sections={sections} />, {
      ...emptySong,
      settings: emptySongSettings
    })

    await user.click(screen.getByRole('button', { name: 'show-details' }))
    expect(screen.queryByRole('button', { name: 'show-details' })).not.toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'hide-details' })).toBeInTheDocument()

    const newSong = { ...emptySong, id: 'new' }
    rerender(
      <SongProvider song={newSong}>
        <SongOutlineProvider>
          <SongOutlineToolbar toggleAdd={vi.fn()} sections={sections} />
        </SongOutlineProvider>
      </SongProvider>
    )
  })
})
