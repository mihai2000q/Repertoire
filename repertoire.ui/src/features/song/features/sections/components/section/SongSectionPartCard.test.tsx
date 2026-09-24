import {
  emptyArtist,
  emptySong,
  emptySongPart,
  reduxRender,
  withToastify
} from '../../../../../../test-utils.tsx'
import SongSectionPartCard from './SongSectionPartCard.tsx'
import { Instrument, SongPart } from '../../../../../../types/models/Song.ts'
import { screen, within } from '@testing-library/react'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import { userEvent } from '@testing-library/user-event'
import { UpdateSongPartRequest } from '../../../parts/types/requests/SongPartRequests.ts'
import { BandMember } from '../../../../../../types/models/Artist.ts'
import { useClickSelect } from '../../../../../../context/ClickSelectContext.tsx'
import { SongProvider } from '../../../../context/SongContext.tsx'
import { SongOutlineProvider } from '../../../outline/context/SongOutlineContext.tsx'
import { ReactNode } from 'react'

vi.mock('../../../../../../context/ClickSelectContext', () => ({
  useClickSelect: vi.fn()
}))

describe('Song Section Part Card', () => {
  const part: SongPart = {
    ...emptySongPart,
    id: '1',
    name: 'Solo 1',
    rehearsals: 12,
    confidence: 50,
    progress: 150
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

  function render(ui: ReactNode, song = emptySong) {
    return reduxRender(
      <SongProvider song={song}>
        <SongOutlineProvider>{ui}</SongOutlineProvider>
      </SongProvider>
    )
  }

  it('should render and display minimal info', () => {
    render(
      <SongSectionPartCard part={part} maxPartProgress={0} isDragging={false} sectionId={''} />
    )

    expect(screen.getByText(part.name)).toBeInTheDocument()
    expect(screen.getByText(part.rehearsals)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'add-rehearsal' })).toBeInTheDocument()
  })

  it('should render and display maximal info', async () => {
    const user = userEvent.setup()

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
    const [{ rerender }] = render(
      <SongSectionPartCard
        part={{
          ...part,
          bandMembers: [bandMember],
          instrument: instrument
        }}
        maxPartProgress={0}
        isDragging={false}
        sectionId={''}
      />,
      { ...emptySong, artist: { ...emptyArtist, isBand: true } }
    )

    expect(screen.getByRole('img', { name: bandMember.name })).toBeInTheDocument()
    expect(screen.getByLabelText('instrument-icon')).toBeInTheDocument()
    expect(screen.getByText(part.name)).toBeInTheDocument()
    expect(screen.getByText(part.rehearsals)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'add-rehearsal' })).toBeInTheDocument()

    await user.hover(screen.getByLabelText('instrument-icon'))
    expect(await screen.findByRole('tooltip', { name: instrument.name })).toBeInTheDocument()

    // when artist is not a band
    const newSong = { ...emptySong, artist: { ...emptyArtist, isBand: false } }
    rerender(
      <SongProvider song={newSong}>
        <SongSectionPartCard
          part={{ ...part, bandMembers: [bandMember], instrument }}
          sectionId={'section-id'}
          maxPartProgress={0}
          isDragging={false}
        />
      </SongProvider>
    )

    expect(screen.queryByRole('img', { name: bandMember.name })).not.toBeInTheDocument()
  })

  it('should show details on click', async () => {
    const user = userEvent.setup()
    const maxPartProgress = 67

    render(
      <SongSectionPartCard
        part={part}
        maxPartProgress={maxPartProgress}
        isDragging={false}
        sectionId={''}
      />
    )

    await user.click(screen.getByLabelText(`song-section-part-${part.name}`))

    expect(await screen.findByLabelText(`song-section-part-details-${part.name}`)).toBeVisible()
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
  })

  it('should display menu on right click', async () => {
    const user = userEvent.setup()

    render(
      <SongSectionPartCard part={part} maxPartProgress={0} isDragging={false} sectionId={''} />
    )

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByLabelText(`song-section-part-${part.name}`)
    })

    expect(within(screen.getByRole('menu')).getByText(/part/i)).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /edit/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should open edit song part modal when clicking edit', async () => {
      const user = userEvent.setup()

      render(
        <SongSectionPartCard part={part} maxPartProgress={0} isDragging={false} sectionId={''} />
      )

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByLabelText(`song-section-part-${part.name}`)
      })
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

      render(
        withToastify(
          <SongSectionPartCard part={part} maxPartProgress={0} isDragging={false} sectionId={''} />
        ),
        { ...emptySong, id: songId, artist: { ...emptyArtist, isBand: false } }
      )

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByLabelText(`song-section-part-${part.name}`)
      })
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

    render(
      <SongSectionPartCard
        part={part}
        maxPartProgress={0}
        isDragging={false}
        showRehearsalsToast={showToast}
        sectionId={''}
      />
    )

    await user.click(screen.getByRole('button', { name: 'add-rehearsal' }))

    expect(capturedRequest).toStrictEqual({
      ...part,
      rehearsals: part.rehearsals + 1
    })
    expect(showToast).toHaveBeenCalledOnce()
  })

  it('should disable context menu and rehearsal button, when click selection is active', async () => {
    const user = userEvent.setup()

    vi.mocked(useClickSelect).mockReturnValue({
      selectables: [],
      addSelectable: vi.fn(),
      removeSelectable: vi.fn(),
      selectedIds: [],
      isClickSelectionActive: true,
      clearSelection: vi.fn()
    })

    render(
      <SongSectionPartCard
        part={part}
        maxPartProgress={0}
        isDragging={false}
        showRehearsalsToast={vi.fn()}
        sectionId={''}
      />
    )

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByLabelText(`song-section-part-${part.name}`)
    })
    expect(screen.queryByRole('menu')).not.toBeInTheDocument()

    expect(screen.getByRole('button', { name: 'add-rehearsal' })).toBeDisabled()
  })

  it('should display a checkmark, when click selected (part of the selected ids)', () => {
    const sectionId = 'section-id'

    vi.mocked(useClickSelect).mockReturnValue({
      selectables: [],
      addSelectable: vi.fn(),
      removeSelectable: vi.fn(),
      selectedIds: [`part-${part.id}:${sectionId}`],
      isClickSelectionActive: true,
      clearSelection: vi.fn()
    })

    render(
      <SongSectionPartCard
        part={part}
        maxPartProgress={0}
        isDragging={false}
        showRehearsalsToast={vi.fn()}
        sectionId={sectionId}
      />
    )
    expect(screen.getByTestId('selected-checkmark')).toBeInTheDocument()
  })

  describe('should be selected', () => {
    it('when hovered', async () => {
      const user = userEvent.setup()

      render(
        <SongSectionPartCard
          part={part}
          maxPartProgress={0}
          isDragging={false}
          showRehearsalsToast={vi.fn()}
          sectionId={''}
        />
      )

      await user.hover(screen.getByLabelText(`song-section-part-${part.name}`))

      expect(screen.getByLabelText(`song-section-part-${part.name}`)).toHaveAttribute(
        'aria-selected',
        'true'
      )
    })

    it('when context menu is open', async () => {
      const user = userEvent.setup()

      render(
        <SongSectionPartCard
          part={part}
          maxPartProgress={0}
          isDragging={false}
          showRehearsalsToast={vi.fn()}
          sectionId={''}
        />
      )

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByLabelText(`song-section-part-${part.name}`)
      })

      expect(screen.getByLabelText(`song-section-part-${part.name}`)).toHaveAttribute(
        'aria-selected',
        'true'
      )
    })

    it('when is dragging', async () => {
      render(
        <SongSectionPartCard
          part={part}
          maxPartProgress={0}
          isDragging={true}
          showRehearsalsToast={vi.fn()}
          sectionId={''}
        />
      )

      expect(screen.getByLabelText(`song-section-part-${part.name}`)).toHaveAttribute(
        'aria-selected',
        'true'
      )
    })

    it('when click selected (part of the selected ids)', () => {
      const sectionId = '123'

      vi.mocked(useClickSelect).mockReturnValue({
        selectables: [],
        addSelectable: vi.fn(),
        removeSelectable: vi.fn(),
        selectedIds: [`part-${part.id}:${sectionId}`],
        isClickSelectionActive: true,
        clearSelection: vi.fn()
      })

      render(
        <SongSectionPartCard
          part={part}
          maxPartProgress={0}
          isDragging={false}
          showRehearsalsToast={vi.fn()}
          sectionId={sectionId}
        />
      )

      expect(screen.getByLabelText(`song-section-part-${part.name}`)).toHaveAttribute(
        'aria-selected',
        'true'
      )
    })
  })
})
