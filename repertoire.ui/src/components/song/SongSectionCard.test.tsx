import { emptySongSection, reduxRender, withToastify } from '../../test-utils.tsx'
import SongSectionCard from './SongSectionCard.tsx'
import { Instrument, SongSection } from './../../types/models/Song.ts'
import { screen } from '@testing-library/react'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import { userEvent } from '@testing-library/user-event'
import { UpdateSongSectionRequest } from '../../types/requests/SongRequests.ts'
import { BandMember } from '../../types/models/Artist.ts'
import { beforeEach, expect } from 'vitest'
import { useClickSelect } from '../../context/ClickSelectContext.tsx'

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
    vi.restoreAllMocks()
    server.resetHandlers()
  })

  beforeAll(() => {
    vi.mock('../../context/ClickSelectContext', () => ({
      useClickSelect: vi.fn()
    }))
    server.listen()
  })

  afterAll(() => server.close())

  it('should render and display minimal info', () => {
    reduxRender(
      <SongSectionCard
        section={section}
        songId={''}
        maxSectionProgress={0}
        maxSectionRehearsals={0}
        showDetails={false}
        isDragging={false}
      />
    )

    expect(screen.getByRole('button', { name: 'drag-handle' })).toBeInTheDocument()
    expect(screen.getByText(section.songSectionType.name)).toBeInTheDocument()
    expect(screen.getByText(section.name)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'add-rehearsal' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()
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
    const [{ rerender }] = reduxRender(
      <SongSectionCard
        section={{
          ...section,
          bandMember: bandMember,
          instrument: instrument
        }}
        songId={''}
        maxSectionProgress={0}
        maxSectionRehearsals={0}
        showDetails={false}
        isDragging={false}
        isArtistBand={true}
      />
    )

    expect(screen.getByRole('button', { name: 'drag-handle' })).toBeInTheDocument()
    expect(screen.getByRole('img', { name: bandMember.name })).toBeInTheDocument()
    expect(screen.getByLabelText('instrument-icon')).toBeInTheDocument()
    expect(screen.getByText(section.name)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'add-rehearsal' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()

    await user.hover(screen.getByRole('img', { name: bandMember.name }))
    expect(await screen.findByText(bandMember.name)).toBeInTheDocument()
    expect(screen.getAllByRole('img', { name: bandMember.name })).toHaveLength(2)
    expect(screen.getByText(bandMember.roles[0].name)).toBeInTheDocument()

    // when artist is not a band
    rerender(
      <SongSectionCard
        section={{
          ...section,
          bandMember: bandMember,
          instrument: instrument
        }}
        songId={''}
        maxSectionProgress={0}
        maxSectionRehearsals={0}
        showDetails={false}
        isDragging={false}
        isArtistBand={false}
      />
    )

    expect(screen.queryByRole('img', { name: bandMember.name })).not.toBeInTheDocument()
  })

  it('should show details', async () => {
    const user = userEvent.setup()
    const maxSectionProgress = 67

    reduxRender(
      <SongSectionCard
        section={section}
        songId={''}
        maxSectionProgress={maxSectionProgress}
        maxSectionRehearsals={0}
        showDetails={true}
        isDragging={false}
      />
    )

    expect(screen.getAllByText(section.rehearsals)).toHaveLength(2) // the one visible and the one in the tooltip
    expect(screen.getByTestId('rehearsals')).toHaveTextContent(section.rehearsals.toString())
    expect(screen.getByRole('progressbar', { name: 'confidence' })).toBeInTheDocument()
    expect(screen.getByRole('progressbar', { name: 'confidence' })).toHaveValue(section.confidence)
    expect(screen.getByRole('progressbar', { name: 'progress' })).toBeInTheDocument()
    expect(screen.getByRole('progressbar', { name: 'progress' })).toHaveValue(
      (section.progress / maxSectionProgress) * 100
    )

    await user.hover(screen.getByTestId('rehearsals'))
    expect(
      screen.getByRole('tooltip', { name: new RegExp(section.rehearsals.toString()) })
    ).toBeInTheDocument()

    await user.hover(screen.getByRole('progressbar', { name: 'confidence' }))
    expect(
      screen.getByRole('tooltip', { name: new RegExp(section.confidence.toString()) })
    ).toBeInTheDocument()

    await user.hover(screen.getByRole('progressbar', { name: 'progress' }))
    expect(
      screen.getByRole('tooltip', { name: new RegExp(section.progress.toString()) })
    ).toBeInTheDocument()
  })

  it('should display menu on right click', async () => {
    const user = userEvent.setup()

    reduxRender(
      <SongSectionCard
        section={section}
        songId={''}
        maxSectionProgress={0}
        maxSectionRehearsals={0}
        showDetails={true}
        isDragging={false}
      />
    )

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByLabelText(`song-section-${section.name}`)
    })

    expect(screen.getByRole('menuitem', { name: /edit/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  it('should display menu by clicking on the dots button', async () => {
    const user = userEvent.setup()

    reduxRender(
      <SongSectionCard
        section={section}
        songId={''}
        maxSectionProgress={0}
        maxSectionRehearsals={0}
        showDetails={true}
        isDragging={false}
      />
    )

    await user.click(screen.getByRole('button', { name: 'more-menu' }))

    expect(screen.getByRole('menuitem', { name: /edit/i })).toBeInTheDocument()
    expect(screen.getByRole('menuitem', { name: /delete/i })).toBeInTheDocument()
  })

  describe('on menu', () => {
    it('should open edit song section modal when clicking edit', async () => {
      const user = userEvent.setup()

      reduxRender(
        <SongSectionCard
          section={section}
          songId={''}
          maxSectionProgress={0}
          maxSectionRehearsals={0}
          showDetails={true}
          isDragging={false}
        />
      )

      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /edit/i }))

      expect(await screen.findByRole('dialog', { name: /edit song section/i })).toBeInTheDocument()
    })

    it('should display warning modal and delete section, when clicking delete', async () => {
      const user = userEvent.setup()

      const songId = 'some-song-id'

      server.use(
        http.delete(`/songs/sections/${section.id}/from/${songId}`, () => {
          return HttpResponse.json({ message: 'it worked' })
        })
      )

      reduxRender(
        withToastify(
          <SongSectionCard
            section={section}
            songId={songId}
            maxSectionProgress={0}
            maxSectionRehearsals={0}
            showDetails={true}
            isDragging={false}
          />
        )
      )

      await user.click(screen.getByRole('button', { name: 'more-menu' }))
      await user.click(screen.getByRole('menuitem', { name: /delete/i }))

      expect(await screen.findByRole('dialog', { name: /delete section/i })).toBeInTheDocument()
      expect(screen.getByRole('heading', { name: /delete section/i })).toBeInTheDocument()
      await user.click(screen.getByRole('button', { name: /yes/i }))

      expect(screen.getByText(`${section.name} deleted!`)).toBeInTheDocument()
    })
  })

  it('should add 1 rehearsal', async () => {
    const user = userEvent.setup()

    let capturedRequest: UpdateSongSectionRequest
    server.use(
      http.put(`/songs/sections`, async (req) => {
        capturedRequest = (await req.request.json()) as UpdateSongSectionRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const showToast = vi.fn()

    reduxRender(
      <SongSectionCard
        section={section}
        songId={''}
        maxSectionProgress={0}
        maxSectionRehearsals={0}
        showDetails={true}
        isDragging={false}
        showRehearsalsToast={showToast}
      />
    )

    await user.click(screen.getByRole('button', { name: 'add-rehearsal' }))

    expect(capturedRequest).toStrictEqual({
      ...section,
      typeId: section.songSectionType.id,
      rehearsals: section.rehearsals + 1
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
      <SongSectionCard
        section={section}
        songId={''}
        maxSectionProgress={0}
        maxSectionRehearsals={0}
        showDetails={true}
        isDragging={false}
        showRehearsalsToast={vi.fn()}
      />
    )

    await user.pointer({
      keys: '[MouseRight>]',
      target: screen.getByLabelText(`song-section-${section.name}`)
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
      selectedIds: [section.id],
      isClickSelectionActive: true,
      clearSelection: vi.fn()
    })

    reduxRender(
      <SongSectionCard
        section={section}
        songId={''}
        maxSectionProgress={0}
        maxSectionRehearsals={0}
        showDetails={false}
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
        <SongSectionCard
          section={section}
          songId={''}
          maxSectionProgress={0}
          maxSectionRehearsals={0}
          showDetails={false}
          isDragging={false}
          showRehearsalsToast={vi.fn()}
        />
      )

      await user.hover(screen.getByLabelText(`song-section-${section.name}`))

      expect(screen.getByLabelText(`song-section-${section.name}`)).toHaveAttribute(
        'aria-selected',
        'true'
      )
    })

    it('when context menu is open', async () => {
      const user = userEvent.setup()

      reduxRender(
        <SongSectionCard
          section={section}
          songId={''}
          maxSectionProgress={0}
          maxSectionRehearsals={0}
          showDetails={false}
          isDragging={false}
          showRehearsalsToast={vi.fn()}
        />
      )

      await user.pointer({
        keys: '[MouseRight>]',
        target: screen.getByLabelText(`song-section-${section.name}`)
      })

      expect(screen.getByLabelText(`song-section-${section.name}`)).toHaveAttribute(
        'aria-selected',
        'true'
      )
    })

    it('when more menu is open', async () => {
      const user = userEvent.setup()

      reduxRender(
        <SongSectionCard
          section={section}
          songId={''}
          maxSectionProgress={0}
          maxSectionRehearsals={0}
          showDetails={false}
          isDragging={false}
          showRehearsalsToast={vi.fn()}
        />
      )

      await user.click(screen.getByRole('button', { name: 'more-menu' }))

      expect(screen.getByLabelText(`song-section-${section.name}`)).toHaveAttribute(
        'aria-selected',
        'true'
      )
    })

    it('when is dragging', async () => {
      reduxRender(
        <SongSectionCard
          section={section}
          songId={''}
          maxSectionProgress={0}
          maxSectionRehearsals={0}
          showDetails={false}
          isDragging={true}
          showRehearsalsToast={vi.fn()}
        />
      )

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
        selectedIds: [section.id],
        isClickSelectionActive: true,
        clearSelection: vi.fn()
      })

      reduxRender(
        <SongSectionCard
          section={section}
          songId={''}
          maxSectionProgress={0}
          maxSectionRehearsals={0}
          showDetails={false}
          isDragging={false}
          showRehearsalsToast={vi.fn()}
        />
      )

      expect(screen.getByLabelText(`song-section-${section.name}`)).toHaveAttribute(
        'aria-selected',
        'true'
      )
    })
  })
})
