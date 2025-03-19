import { emptySongSettings, reduxRender, withToastify } from '../../test-utils.tsx'
import AddNewSongSection from './AddNewSongSection.tsx'
import { screen } from '@testing-library/react'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { userEvent } from '@testing-library/user-event'
import { Instrument, SongSectionType } from '../../types/models/Song.ts'
import { CreateSongSectionRequest } from '../../types/requests/SongRequests.ts'
import { BandMember } from '../../types/models/Artist.ts'

describe('Add New Song Section', () => {
  const sectionTypes: SongSectionType[] = [
    {
      id: '1',
      name: 'Solo'
    },
    {
      id: '2',
      name: 'Riff'
    }
  ]

  const instruments: Instrument[] = [
    {
      id: '1',
      name: 'Guitar'
    },
    {
      id: '2',
      name: 'Piano'
    },
    {
      id: '3',
      name: 'Flute'
    }
  ]

  const bandMembers: BandMember[] = [
    {
      id: '1',
      name: 'Chester',
      roles: []
    },
    {
      id: '2',
      name: 'Michael',
      roles: []
    }
  ]

  const handlers = [
    http.get('/songs/instruments', async () => {
      return HttpResponse.json(instruments)
    }),
    http.get('/songs/sections/types', async () => {
      return HttpResponse.json(sectionTypes)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    const [{ rerender }] = reduxRender(
      <AddNewSongSection
        opened={true}
        onClose={() => {}}
        songId={''}
        settings={emptySongSettings}
      />
    )

    expect(screen.getByRole('button', { name: 'select-band-member' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'select-band-member' })).toBeDisabled()
    expect(screen.getByRole('button', { name: 'select-instrument' })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /song-section-type/i })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /name/i })).toBeInTheDocument()
    expect(await screen.findByRole('textbox', { name: /name/i })).toHaveFocus()
    expect(screen.getByRole('button', { name: /add/i })).toBeInTheDocument()

    expect(screen.getByRole('textbox', { name: /song-section-type/i })).not.toBeInvalid()
    expect(screen.getByRole('textbox', { name: /name/i })).not.toBeInvalid()

    rerender(
      <AddNewSongSection
        opened={true}
        onClose={() => {}}
        songId={''}
        settings={emptySongSettings}
        bandMembers={[]}
      />
    )
    expect(screen.getByRole('button', { name: 'select-band-member' })).not.toBeDisabled()
  })

  it('should have default options based on settings', async () => {
    const defaultInstrument = instruments[1]
    const defaultBandMember = bandMembers[1]

    reduxRender(
      <AddNewSongSection
        opened={true}
        onClose={() => {}}
        songId={''}
        settings={{ ...emptySongSettings, defaultBandMember, defaultInstrument }}
        bandMembers={bandMembers}
      />
    )

    expect(screen.getByRole('button', { name: defaultBandMember.name })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: defaultInstrument.name })).toBeInTheDocument()
  })

  it('should send create request when name is typed and type is selected', async () => {
    const user = userEvent.setup()

    const onClose = vitest.fn()
    const songId = 'some id'

    const newSectionType = sectionTypes[0]
    const newName = 'Section 1'

    let capturedRequest: CreateSongSectionRequest
    server.use(
      http.post('/songs/sections', async (req) => {
        capturedRequest = (await req.request.json()) as CreateSongSectionRequest
        return HttpResponse.json({ message: 'section added!' })
      })
    )

    reduxRender(
      withToastify(
        <AddNewSongSection
          opened={true}
          onClose={onClose}
          songId={songId}
          settings={emptySongSettings}
        />
      )
    )

    await user.click(screen.getByRole('textbox', { name: /song-section-type/i }))
    await user.click(await screen.findByText(newSectionType.name))

    await user.type(screen.getByRole('textbox', { name: /name/i }), newName)

    await user.click(screen.getByRole('button', { name: /add/i }))

    expect(capturedRequest).toStrictEqual({
      typeId: newSectionType.id,
      name: newName,
      songId: songId
    })
    expect(onClose).toHaveBeenCalledOnce()
    expect(screen.getByRole('textbox', { name: /name/i })).toHaveValue('')
    expect(await screen.findByRole('textbox', { name: /song-section-type/i })).toHaveValue('')

    expect(screen.getByText(`${newName} added!`)).toBeInTheDocument()
  })

  it('should send create request when all fields are filled', async () => {
    const user = userEvent.setup()

    const onClose = vitest.fn()
    const songId = 'some id'

    const newSectionType = sectionTypes[0]
    const newName = 'Section 1'
    const newInstrument = instruments[0]
    const newBandMember = bandMembers[0]

    let capturedRequest: CreateSongSectionRequest
    server.use(
      http.post('/songs/sections', async (req) => {
        capturedRequest = (await req.request.json()) as CreateSongSectionRequest
        return HttpResponse.json({ message: 'section added!' })
      })
    )

    reduxRender(
      withToastify(
        <AddNewSongSection
          opened={true}
          onClose={onClose}
          songId={songId}
          bandMembers={bandMembers}
          settings={emptySongSettings}
        />
      )
    )

    // fill fields
    await user.click(screen.getByRole('button', { name: 'select-band-member' }))
    await user.click(await screen.findByRole('option', { name: newBandMember.name }))

    await user.click(screen.getByRole('button', { name: 'select-instrument' }))
    await user.click(await screen.findByRole('option', { name: newInstrument.name }))

    await user.click(screen.getByRole('textbox', { name: /song-section-type/i }))
    await user.click(await screen.findByText(newSectionType.name))

    await user.type(screen.getByRole('textbox', { name: /name/i }), newName)

    await user.click(screen.getByRole('button', { name: /add/i }))

    expect(capturedRequest).toStrictEqual({
      bandMemberId: newBandMember.id,
      instrumentId: newInstrument.id,
      typeId: newSectionType.id,
      name: newName,
      songId: songId
    })
    expect(onClose).toHaveBeenCalledOnce()

    // reset fields
    expect(screen.getByRole('button', { name: 'select-band-member' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'select-instrument' })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /name/i })).toHaveValue('')
    expect(screen.getByRole('textbox', { name: /song-section-type/i })).toHaveValue('')

    expect(screen.getByText(`${newName} added!`)).toBeInTheDocument()
  })

  // Validation

  it('should display error when name was typed and then removed', async () => {
    const user = userEvent.setup()

    const newName = 'New Name'

    reduxRender(
      <AddNewSongSection
        opened={true}
        onClose={() => {}}
        songId={''}
        settings={emptySongSettings}
      />
    )

    await user.type(screen.getByRole('textbox', { name: /name/i }), newName)
    await user.clear(screen.getByRole('textbox', { name: /name/i }))

    expect(screen.getByRole('textbox', { name: /name/i })).toBeInvalid()
  })

  it('should not send create request when neither type is selected nor name is typed', async () => {
    const user = userEvent.setup()

    let capturedRequest: CreateSongSectionRequest
    server.use(
      http.post('/songs/sections', async (req) => {
        capturedRequest = (await req.request.json()) as CreateSongSectionRequest
        return HttpResponse.json({ message: 'section added!' })
      })
    )

    reduxRender(
      <AddNewSongSection
        opened={true}
        onClose={() => {}}
        songId={''}
        settings={emptySongSettings}
      />
    )

    await user.click(screen.getByRole('button', { name: /add/i }))

    expect(screen.getByRole('textbox', { name: /song-section-type/i })).toBeInvalid()
    expect(screen.getByRole('textbox', { name: /name/i })).toBeInvalid()

    expect(capturedRequest).toBeUndefined()
  })

  it('should not send create request when name is not typed', async () => {
    const user = userEvent.setup()

    const newSectionType = sectionTypes[0]

    let capturedRequest: CreateSongSectionRequest
    server.use(
      http.post('/songs/sections', async (req) => {
        capturedRequest = (await req.request.json()) as CreateSongSectionRequest
        return HttpResponse.json({ message: 'section added!' })
      })
    )

    reduxRender(
      <AddNewSongSection
        opened={true}
        onClose={() => {}}
        songId={''}
        settings={emptySongSettings}
      />
    )

    await user.click(screen.getByRole('textbox', { name: /song-section-type/i }))
    await user.click(await screen.findByText(newSectionType.name))

    await user.click(screen.getByRole('button', { name: /add/i }))

    expect(screen.getByRole('textbox', { name: /name/i })).toBeInvalid()
    expect(screen.getByRole('textbox', { name: /song-section-type/i })).not.toBeInvalid()

    expect(capturedRequest).toBeUndefined()
  })

  it('should not send create request when type is not selected', async () => {
    const user = userEvent.setup()

    const newName = 'Section 1'

    let capturedRequest: CreateSongSectionRequest
    server.use(
      http.post('/songs/sections', async (req) => {
        capturedRequest = (await req.request.json()) as CreateSongSectionRequest
        return HttpResponse.json({ message: 'section added!' })
      })
    )

    reduxRender(
      <AddNewSongSection
        opened={true}
        onClose={() => {}}
        songId={''}
        settings={emptySongSettings}
      />
    )

    await user.type(screen.getByRole('textbox', { name: /name/i }), newName)

    await user.click(screen.getByRole('button', { name: /add/i }))

    expect(screen.getByRole('textbox', { name: /name/i })).not.toBeInvalid()
    expect(screen.getByRole('textbox', { name: /song-section-type/i })).toBeInvalid()

    expect(capturedRequest).toBeUndefined()
  })

  it('should refresh errors when reopened', async () => {
    const user = userEvent.setup()

    const uut = (opened = true) => (
      <AddNewSongSection
        opened={opened}
        onClose={() => {}}
        songId={''}
        settings={emptySongSettings}
      />
    )

    const [{ rerender }] = reduxRender(uut())

    await user.click(screen.getByRole('button', { name: /add/i }))

    expect(screen.getByRole('textbox', { name: /song-section-type/i })).toBeInvalid()
    expect(screen.getByRole('textbox', { name: /name/i })).toBeInvalid()

    rerender(uut(false))
    rerender(uut())

    expect(screen.getByRole('textbox', { name: /song-section-type/i })).not.toBeInvalid()
    expect(screen.getByRole('textbox', { name: /name/i })).not.toBeInvalid()
  })
})
