import {
  emptyArtist,
  emptySong,
  emptySongPart,
  emptySongSection,
  reduxRender,
  withToastify
} from '../../../../../../../test-utils.tsx'
import SongPartCard from './SongPartCard.tsx'
import { Instrument, SongPart } from '../../../../../../../types/models/Song.ts'
import { act, screen } from '@testing-library/react'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import { userEvent } from '@testing-library/user-event'
import { UpdateSongPartRequest } from '../types/requests/SongPartRequests.ts'
import { BandMember } from '../../../../../../../types/models/Artist.ts'
import { useClickSelect } from '../../../../../../../context/ClickSelectContext.tsx'
import OutlineView from '../../outline/types/enums/OutlineView.ts'
import { setSong } from '../../../state/slice/songSlice.tsx'

// Mock Context
vi.mock('../../../../../../../context/ClickSelectContext', () => ({
  useClickSelect: vi.fn()
}))

describe('Song Part Card', () => {
  const part: SongPart = {
    ...emptySongPart,
    id: '1',
    name: 'Solo 1',
    rehearsals: 12,
    confidence: 50,
    progress: 150,
    sections: [{ ...emptySongSection, id: '1', name: 'Solo' }]
  }

  const handlers = [
    http.get(`/songs/instruments`, () => {
      return HttpResponse.json([])
    })
  ]

  const server = setupServer(...handlers)

  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(useClickSelect).mockReturnValue({
      selectables: [],
      addSelectable: vi.fn(),
      removeSelectable: vi.fn(),
      selectedIds: [],
      isClickSelectionActive: false,
      clearSelection: vi.fn()
    })
  })

  afterEach(() => {
    server.resetHandlers()
    vi.restoreAllMocks()
  })

  beforeAll(() => server.listen())

  afterAll(() => {
    vi.restoreAllMocks()
    server.close()
  })

  it('should render and display minimal info', () => {
    reduxRender(<SongPartCard part={part} maxPartProgress={0} isDragging={false} />)

    expect(screen.getByRole('button', { name: 'drag-handle' })).toBeInTheDocument()
    expect(screen.getByText(part.name)).toBeInTheDocument()
    expect(screen.getByText(part.rehearsals)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'add-rehearsal' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()
  })

  it('should render and display maximal info', async () => {
    const bandMember: BandMember = {
      id: '1',
      name: 'Mike',
      roles: [{ id: '1', name: 'Guitarist' }],
      imageUrl: 'default.png'
    }

    const instrument: Instrument = {
      id: '1',
      name: 'Electric Guitar'
    }

    // when artist is a band
    const [_, store] = reduxRender(
      <SongPartCard
        part={{
          ...part,
          bandMember: bandMember,
          instrument: instrument
        }}
        maxPartProgress={0}
        isDragging={false}
      />,
      {
        song: { songId: '', isArtistBand: true }
      }
    )

    expect(screen.getByRole('button', { name: 'drag-handle' })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: bandMember.name })).toBeInTheDocument()
    expect(screen.getByLabelText('instrument-icon')).toBeInTheDocument()
    expect(screen.getByText(part.name)).toBeInTheDocument()
    expect(screen.getByText(part.rehearsals)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'add-rehearsal' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()

    // when artist is not a band
    const newSong = { ...emptySong, artist: { ...emptyArtist, isBand: false } }
    await act(() => store.dispatch(setSong(newSong)))

    expect(screen.queryByRole('img', { name: bandMember.name })).not.toBeInTheDocument()
  })

  async function shouldShowDetails(maxPartProgress: number = 0) {
    const user = userEvent.setup()

    expect(await screen.findByLabelText(`song-part-details-${part.name}`)).toBeVisible()
    expect(screen.getByRole('progressbar', { name: 'confidence' })).toBeInTheDocument()
    expect(screen.getByRole('progressbar', { name: 'confidence' })).toHaveValue(part.confidence)
    expect(screen.getByRole('progressbar', { name: 'progress' })).toBeInTheDocument()
    expect(screen.getByRole('progressbar', { name: 'progress' })).toHaveValue(
      (part.progress / maxPartProgress) * 100
    )

    await user.hover(screen.getByRole('progressbar', { name: 'confidence' }))
    expect(
      screen.getByRole('tooltip', { name: new RegExp(part.confidence.toString()) })
    ).toBeInTheDocument()

    await user.hover(screen.getByRole('progressbar', { name: 'progress' }))
    expect(
      screen.getByRole('tooltip', { name: new RegExp(part.progress.toString()) })
    ).toBeInTheDocument()
  }

  it('should show details on click', async () => {
    const user = userEvent.setup()
    const maxPartProgress = 67

    reduxRender(<SongPartCard part={part} maxPartProgress={maxPartProgress} isDragging={false} />)

    await user.click(screen.getByLabelText(`song-part-${part.name}`))

    await shouldShowDetails(maxPartProgress)
  })

  it('should show details from redux selector', async () => {
    const maxPartProgress = 67

    reduxRender(<SongPartCard part={part} maxPartProgress={maxPartProgress} isDragging={false} />, {
      song: { songId: '', isArtistBand: false },
      songOutline: { view: OutlineView.Parts, showDetails: true }
    })

    await shouldShowDetails(maxPartProgress)
  })

  it('should display menu on right click', async () => {
    const user = userEvent.setup()

    reduxRender(<SongPartCard part={part} maxPartProgress={0} isDragging={false} />)

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByLabelText(`song-part-${part.name}`)
    })

    expect(screen.getByRole('menuitem', { name: /edit/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  it('should display menu by clicking on the dots button', async () => {
    const user = userEvent.setup()

    reduxRender(<SongPartCard part={part} maxPartProgress={0} isDragging={false} />)

    await user.click(screen.getByRole('button', { name: 'more-menu' }))

    expect(screen.getByRole('menuitem', { name: /edit/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should open edit song part modal when clicking edit', async () => {
      const user = userEvent.setup()

      reduxRender(<SongPartCard part={part} maxPartProgress={0} isDragging={false} />)

      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /edit/i }))

      expect(await screen.findByRole('dialog', { name: /edit song part/i })).toBeInTheDocument()
    })

    it('should display warning modal and delete part, when clicking delete', async () => {
      const user = userEvent.setup()

      const songId = 'some-song-id'

      server.use(
        http.delete(`/songs/parts/${part.id}/from/${songId}`, () => {
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      reduxRender(
        withToastify(<SongPartCard part={part} maxPartProgress={0} isDragging={false} />),
        { song: { songId: songId, isArtistBand: false } }
      )

      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /delete/i }))

      expect(await screen.findByRole('dialog', { name: /delete part/i })).toBeInTheDocument()
      expect(screen.getByRole('heading', { name: /delete part/i })).toBeInTheDocument()
      await user.click(screen.getByRole('button', { name: /yes/i }))

      expect(screen.getByText(`${part.name} deleted!`)).toBeInTheDocument()
    })
  })

  it('should add 1 rehearsal', async () => {
    const user = userEvent.setup()

    let capturedRequest: UpdateSongPartRequest
    server.use(
      http.put(`/songs/parts`, async (req) => {
        capturedRequest = (await req.request.json()) as UpdateSongPartRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const showToast = vi.fn()

    reduxRender(
      <SongPartCard
        part={part}
        maxPartProgress={0}
        isDragging={false}
        showRehearsalsToast={showToast}
      />
    )

    await user.click(screen.getByRole('button', { name: 'add-rehearsal' }))

    expect(capturedRequest).toStrictEqual({
      ...part,
      sectionIds: part.sections.map((section) => section.id),
      rehearsals: part.rehearsals + 1
    })
    expect(showToast).toHaveBeenCalledOnce()
  })

  it('should disable context menu; drag handle, more menu and rehearsal buttons, when click selection is active', async () => {
    const user = userEvent.setup()

    vi.mocked(useClickSelect).mockReturnValue({
      selectables: [],
      addSelectable: vi.fn(),
      removeSelectable: vi.fn(),
      selectedIds: [],
      isClickSelectionActive: true,
      clearSelection: vi.fn()
    })

    reduxRender(
      <SongPartCard
        part={part}
        maxPartProgress={0}
        isDragging={false}
        showRehearsalsToast={vi.fn()}
      />
    )

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByLabelText(`song-part-${part.name}`)
    })
    expect(screen.queryByRole('menu')).not.toBeInTheDocument()

    expect(screen.getByRole('button', { name: 'drag-handle' })).toBeDisabled()
    expect(screen.getByRole('button', { name: 'add-rehearsal' })).toBeDisabled()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeDisabled()
  })

  it('should hide the drag handle and display a checkmark, when click selected (part of the selected ids)', () => {
    vi.mocked(useClickSelect).mockReturnValue({
      selectables: [],
      addSelectable: vi.fn(),
      removeSelectable: vi.fn(),
      selectedIds: [part.id],
      isClickSelectionActive: true,
      clearSelection: vi.fn()
    })

    reduxRender(
      <SongPartCard
        part={part}
        maxPartProgress={0}
        isDragging={false}
        showRehearsalsToast={vi.fn()}
      />
    )

    expect(screen.queryByRole('button', { name: 'drag-handle' })).not.toBeInTheDocument()
    expect(screen.getByTestId('selected-checkmark')).toBeInTheDocument()
  })

  describe('should be selected', () => {
    it('when avatar is hovered', async () => {
      const user = userEvent.setup()

      reduxRender(
        <SongPartCard
          part={part}
          maxPartProgress={0}
          isDragging={false}
          showRehearsalsToast={vi.fn()}
        />
      )

      await user.hover(screen.getByLabelText(`song-part-${part.name}`))

      expect(screen.getByLabelText(`song-part-${part.name}`)).toHaveAttribute(
        'aria-selected',
        'true'
      )
    })

    it('when context menu is open', async () => {
      const user = userEvent.setup()

      reduxRender(
        <SongPartCard
          part={part}
          maxPartProgress={0}
          isDragging={false}
          showRehearsalsToast={vi.fn()}
        />
      )

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByLabelText(`song-part-${part.name}`)
      })

      expect(screen.getByLabelText(`song-part-${part.name}`)).toHaveAttribute(
        'aria-selected',
        'true'
      )
    })

    it('when more menu is open', async () => {
      const user = userEvent.setup()

      reduxRender(
        <SongPartCard
          part={part}
          maxPartProgress={0}
          isDragging={false}
          showRehearsalsToast={vi.fn()}
        />
      )

      await user.click(screen.getByRole('button', { name: 'more-menu' }))

      expect(screen.getByLabelText(`song-part-${part.name}`)).toHaveAttribute(
        'aria-selected',
        'true'
      )
    })

    it('when is dragging', async () => {
      reduxRender(
        <SongPartCard
          part={part}
          maxPartProgress={0}
          isDragging={true}
          showRehearsalsToast={vi.fn()}
        />
      )

      expect(screen.getByLabelText(`song-part-${part.name}`)).toHaveAttribute(
        'aria-selected',
        'true'
      )
    })

    it('when click selected (part of the selected ids)', () => {
      vi.mocked(useClickSelect).mockReturnValue({
        selectables: [],
        addSelectable: vi.fn(),
        removeSelectable: vi.fn(),
        selectedIds: [part.id],
        isClickSelectionActive: true,
        clearSelection: vi.fn()
      })

      reduxRender(
        <SongPartCard
          part={part}
          maxPartProgress={0}
          isDragging={false}
          showRehearsalsToast={vi.fn()}
        />
      )

      expect(screen.getByLabelText(`song-part-${part.name}`)).toHaveAttribute(
        'aria-selected',
        'true'
      )
    })
  })
})
