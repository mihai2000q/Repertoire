import { reduxRender, withToastify } from '../../test-utils.tsx'
import SongSection from './SongSection.tsx'
import { SongSection as SongSectionModel } from './../../types/models/Song.ts'
import { screen } from '@testing-library/react'
import { setupServer } from 'msw/node'
import { http, HttpResponse } from 'msw'
import { userEvent } from '@testing-library/user-event'
import { UpdateSongSectionRequest } from '../../types/requests/SongRequests.ts'

describe('Song Section', () => {
  const section: SongSectionModel = {
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
  const draggableProvided = {
    draggableProps: undefined,
    dragHandleProps: undefined,
    innerRef: undefined
  }

  const handlers = [
    http.get(`/songs/sections/types`, async () => {
      return HttpResponse.json([])
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', () => {
    reduxRender(
      <SongSection
        section={section}
        songId={''}
        maxSectionProgress={0}
        showDetails={false}
        isDragging={false}
        draggableProvided={draggableProvided}
      />
    )

    expect(screen.getByRole('button', { name: 'drag-handle' })).toBeInTheDocument()
    expect(screen.getByText(section.songSectionType.name)).toBeInTheDocument()
    expect(screen.getByText(section.name)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'add-rehearsal' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'more-menu' })).toBeInTheDocument()
  })

  it('should show details', async () => {
    const user = userEvent.setup()
    const maxSectionProgress = 67

    reduxRender(
      <SongSection
        section={section}
        songId={''}
        maxSectionProgress={maxSectionProgress}
        showDetails={true}
        isDragging={false}
        draggableProvided={draggableProvided}
      />
    )

    expect(screen.getByText(section.rehearsals)).toBeInTheDocument()
    expect(screen.getByRole('progressbar', { name: 'confidence' })).toBeInTheDocument()
    expect(screen.getByRole('progressbar', { name: 'confidence' })).toHaveValue(section.confidence)
    expect(screen.getByRole('progressbar', { name: 'progress' })).toBeInTheDocument()
    expect(screen.getByRole('progressbar', { name: 'progress' })).toHaveValue(
      (section.progress / maxSectionProgress) * 100
    )

    await user.hover(screen.getByText(section.rehearsals))
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

  it('should display menu', async () => {
    const user = userEvent.setup()

    reduxRender(
      <SongSection
        section={section}
        songId={''}
        maxSectionProgress={0}
        showDetails={true}
        isDragging={false}
        draggableProvided={draggableProvided}
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
        <SongSection
          section={section}
          songId={''}
          maxSectionProgress={0}
          showDetails={true}
          isDragging={false}
          draggableProvided={draggableProvided}
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
          <SongSection
            section={section}
            songId={songId}
            maxSectionProgress={0}
            showDetails={true}
            isDragging={false}
            draggableProvided={draggableProvided}
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

    reduxRender(
      withToastify(
        <SongSection
          section={section}
          songId={''}
          maxSectionProgress={0}
          showDetails={true}
          isDragging={false}
          draggableProvided={draggableProvided}
        />
      )
    )

    await user.click(screen.getByRole('button', { name: 'add-rehearsal' }))

    expect(capturedRequest).toStrictEqual({
      ...section,
      typeId: section.songSectionType.id,
      rehearsals: section.rehearsals + 1
    })

    expect(
      screen.getByText(new RegExp(`${section.name} rehearsals' .* increased`))
    ).toBeInTheDocument()
  })
})
