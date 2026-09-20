import {
  emptyArtist,
  emptySong,
  emptySongPart,
  emptySongSection,
  reduxRender,
  withToastify
} from '../../../../../test-utils.tsx'
import SongSectionCard from './SongSectionCard.tsx'
import { Instrument, SongSection } from '../../../../../types/models/Song.ts'
import { screen, within } from '@testing-library/react'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import { userEvent } from '@testing-library/user-event'
import { BandMember } from '../../../../../types/models/Artist.ts'
import { useClickSelect } from '../../../../../context/ClickSelectContext.tsx'
import { BulkUpdateSongPartsRequest } from '../../parts/types/requests/SongPartRequests.ts'
import { SongProvider } from '../../../context/SongContext.tsx'
import { SongOutlineProvider } from '../../outline/context/SongOutlineContext.tsx'
import { ReactNode } from 'react'

// Mock Context
vi.mock('../../../../../context/ClickSelectContext', () => ({
  useClickSelect: vi.fn()
}))

describe('Song Section Card', () => {
  const section: SongSection = {
    ...emptySongSection,
    id: '1',
    name: 'Solo 1',
    rehearsals: 12,
    confidence: 50,
    progress: 150,
    songSectionType: {
      id: 'some id',
      name: 'Solo'
    }
  }

  const handlers = [
    http.get(`/songs/sections/types`, () => {
      return HttpResponse.json([])
    }),
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

  function render(
    ui: ReactNode,
    song = emptySong,
    initialDetailsSectionIds: Set<string> = new Set()
  ) {
    return reduxRender(
      <SongProvider song={song}>
        <SongOutlineProvider initialDetailsSectionIds={initialDetailsSectionIds}>
          {ui}
        </SongOutlineProvider>
      </SongProvider>
    )
  }

  it('should render and display info', () => {
    const maxSectionProgress = 60

    render(
      <SongSectionCard
        section={section}
        maxSectionProgress={maxSectionProgress}
        isDragging={false}
      />
    )

    expect(screen.getByText(section.name)).toBeInTheDocument()
    expect(screen.getByText(section.songSectionType.name)).toBeInTheDocument()
    expect(screen.getByText(section.rehearsals)).toBeInTheDocument()
    expect(screen.getByRole('progressbar', { name: 'confidence' })).toBeInTheDocument()
    expect(screen.getByRole('progressbar', { name: 'confidence' })).toHaveValue(section.confidence)
    expect(screen.getByRole('progressbar', { name: 'progress' })).toBeInTheDocument()
    expect(screen.getByRole('progressbar', { name: 'progress' })).toHaveValue(
      (section.progress / maxSectionProgress) * 100
    )
  })

  it('should render and display aggregated band members and instruments', async () => {
    const user = userEvent.setup()
    const bandMembers: BandMember[] = [
      {
        id: '1',
        name: 'Mike',
        roles: [{ id: '1', name: 'Guitarist' }],
        imageUrl: 'default.png'
      },
      {
        id: '2',
        name: 'Leonard',
        roles: [{ id: '2', name: 'Voice' }]
      }
    ]

    const instruments: Instrument[] = [
      {
        id: '1',
        name: 'Electric Guitar'
      },
      {
        id: '2',
        name: 'Voice'
      }
    ]

    // when artist is a band
    const [{ rerender }] = render(
      <SongSectionCard
        section={{
          ...section,
          parts: [
            {
              ...emptySongPart,
              id: '1',
              instrument: instruments[0],
              bandMembers: [bandMembers[0], bandMembers[1]]
            },
            {
              ...emptySongPart,
              id: '2',
              instrument: instruments[1],
              bandMembers: [bandMembers[1]]
            },
            {
              ...emptySongPart,
              id: '3',
              instrument: instruments[0],
              bandMembers: [bandMembers[1]]
            },
            { ...emptySongPart, id: '4', bandMembers: [bandMembers[0]] },
            { ...emptySongPart, id: '5' }
          ]
        }}
        maxSectionProgress={0}
        isDragging={false}
      />,
      { ...emptySong, artist: { ...emptyArtist, isBand: true } }
    )

    const instrumentsEl = screen.getByLabelText('instruments')
    expect(instrumentsEl).toBeInTheDocument()
    instruments.forEach((instrument) => {
      expect(within(instrumentsEl).getByLabelText(instrument.name)).toBeInTheDocument()
    })
    await user.hover(instrumentsEl)
    for (const instrument of instruments) {
      expect(await screen.findByRole('tooltip')).toHaveTextContent(instrument.name)
    }

    const bandMembersEl = screen.getByLabelText('band-members')
    expect(bandMembersEl).toBeInTheDocument()
    bandMembers.forEach((bandMember) => {
      if (bandMember.imageUrl) {
        expect(
          within(bandMembersEl).getByRole('img', { name: bandMember.name })
        ).toBeInTheDocument()
      } else {
        expect(
          within(bandMembersEl).getByLabelText(`default-icon-${bandMember.name}`)
        ).toBeInTheDocument()
      }
    })
    await user.hover(bandMembersEl)
    for (const bandMember of bandMembers) {
      expect(await screen.findByRole('tooltip')).toHaveTextContent(bandMember.name)
    }

    // when artist is not a band
    const newSong = { ...emptySong, artist: { ...emptyArtist, isBand: false } }
    rerender(
      <SongProvider song={newSong}>
        <SongOutlineProvider>
          <SongSectionCard section={section} maxSectionProgress={0} isDragging={false} />
        </SongOutlineProvider>
      </SongProvider>
    )

    expect(screen.queryByLabelText('band-members')).not.toBeInTheDocument()
  })

  it('should show details when clicking', async () => {
    const user = userEvent.setup()

    render(<SongSectionCard section={section} maxSectionProgress={0} isDragging={false} />)

    await user.click(screen.getByLabelText(`song-section-${section.name}`))

    section.parts.forEach((part) => {
      expect(screen.getByLabelText(`song-section-part-${part.name}`)).toBeInTheDocument()
    })
    expect(
      screen.getByLabelText(`add-new-song-section-part-card-${section.name}`)
    ).toBeInTheDocument()
  })

  it('should show details from context', async () => {
    render(
      <SongSectionCard section={section} maxSectionProgress={0} isDragging={false} />,
      emptySong,
      new Set([section.id])
    )

    section.parts.forEach((part) => {
      expect(screen.getByLabelText(`song-section-part-${part.name}`)).toBeInTheDocument()
    })
    expect(
      screen.getByLabelText(`add-new-song-section-part-card-${section.name}`)
    ).toBeInTheDocument()
  })

  it('should display menu on right click', async () => {
    const user = userEvent.setup()

    render(<SongSectionCard section={section} maxSectionProgress={0} isDragging={false} />)

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByLabelText(`song-section-${section.name}`)
    })

    expect(screen.getByRole('menuitem', { name: /edit/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /add rehearsal/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should open edit song section modal when clicking edit', async () => {
      const user = userEvent.setup()

      render(<SongSectionCard section={section} maxSectionProgress={0} isDragging={false} />)

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByLabelText(`song-section-${section.name}`)
      })
      await user.click(screen.getByRole('menuitem', { name: /edit/i }))

      expect(await screen.findByRole('dialog', { name: /edit song section/i })).toBeInTheDocument()
    })

    it('should open add rehearsals to parts when clicking add rehearsal', async () => {
      const user = userEvent.setup()

      let capturedRequest: BulkUpdateSongPartsRequest
      server.use(
        http.put(`/songs/parts/bulk-update`, async (req) => {
          capturedRequest = (await req.request.json()) as BulkUpdateSongPartsRequest
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      const songId = '1'
      render(
        withToastify(
          <SongSectionCard section={section} maxSectionProgress={0} isDragging={false} />
        ),
        { ...emptySong, id: songId, artist: { ...emptyArtist, isBand: false } }
      )

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByLabelText(`song-section-${section.name}`)
      })
      await user.click(screen.getByRole('menuitem', { name: /add rehearsal/i }))
      await user.click(screen.getByRole('button', { name: /confirm/i })) // menu item confirmation

      expect(
        screen.getByText(`Rehearsals added to ${section.name}'s ${section.parts.length} parts!`)
      ).toBeInTheDocument()
      expect(capturedRequest).toStrictEqual({
        requests: section.parts.map((p) => ({
          id: p.id,
          rehearsals: p.rehearsals + 1,
          confidence: p.confidence
        })),
        songId: songId
      })
    })

    it('should open edit song section modal when clicking edit', async () => {
      const user = userEvent.setup()

      render(<SongSectionCard section={section} maxSectionProgress={0} isDragging={false} />)

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByLabelText(`song-section-${section.name}`)
      })
      await user.click(screen.getByRole('menuitem', { name: /edit/i }))

      expect(await screen.findByRole('dialog', { name: /edit song section/i })).toBeInTheDocument()
    })

    it('should display delete section modal when clicking delete', async () => {
      const user = userEvent.setup()

      render(
        withToastify(
          <SongSectionCard section={section} maxSectionProgress={0} isDragging={false} />
        )
      )

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByLabelText(`song-section-${section.name}`)
      })
      await user.click(screen.getByRole('menuitem', { name: /delete/i }))

      expect(await screen.findByRole('dialog', { name: /delete section/i })).toBeInTheDocument()
    })
  })

  it('should disable context menu, when click selection is active', async () => {
    const user = userEvent.setup()

    vi.mocked(useClickSelect).mockReturnValue({
      selectables: [],
      addSelectable: vi.fn(),
      removeSelectable: vi.fn(),
      selectedIds: [],
      isClickSelectionActive: true,
      clearSelection: vi.fn()
    })

    render(<SongSectionCard section={section} maxSectionProgress={0} isDragging={false} />)

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByLabelText(`song-section-${section.name}`)
    })
    expect(screen.queryByRole('menu')).not.toBeInTheDocument()
  })

  it('should display a checkmark, when click selected (part of the selected ids)', () => {
    vi.mocked(useClickSelect).mockReturnValue({
      selectables: [],
      addSelectable: vi.fn(),
      removeSelectable: vi.fn(),
      selectedIds: ['section-' + section.id],
      isClickSelectionActive: true,
      clearSelection: vi.fn()
    })

    render(<SongSectionCard section={section} maxSectionProgress={0} isDragging={false} />)

    expect(screen.getByTestId('selected-checkmark')).toBeInTheDocument()
  })

  describe('should be selected', () => {
    it('when avatar is hovered', async () => {
      const user = userEvent.setup()

      render(<SongSectionCard section={section} maxSectionProgress={0} isDragging={false} />)

      await user.hover(screen.getByLabelText(`song-section-${section.name}`))

      expect(screen.getByLabelText(`song-section-${section.name}`)).toHaveAttribute(
        'aria-selected',
        'true'
      )
    })

    it('when context menu is open', async () => {
      const user = userEvent.setup()

      render(<SongSectionCard section={section} maxSectionProgress={0} isDragging={false} />)

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByLabelText(`song-section-${section.name}`)
      })

      expect(screen.getByLabelText(`song-section-${section.name}`)).toHaveAttribute(
        'aria-selected',
        'true'
      )
    })

    it('when is dragging', async () => {
      render(<SongSectionCard section={section} maxSectionProgress={0} isDragging={true} />)

      expect(screen.getByLabelText(`song-section-${section.name}`)).toHaveAttribute(
        'aria-selected',
        'true'
      )
    })

    it('when click selected (part of the selected ids)', () => {
      vi.mocked(useClickSelect).mockReturnValue({
        selectables: [],
        addSelectable: vi.fn(),
        removeSelectable: vi.fn(),
        selectedIds: ['section-' + section.id],
        isClickSelectionActive: true,
        clearSelection: vi.fn()
      })

      render(<SongSectionCard section={section} maxSectionProgress={0} isDragging={false} />)

      expect(screen.getByLabelText(`song-section-${section.name}`)).toHaveAttribute(
        'aria-selected',
        'true'
      )
    })
  })
})
