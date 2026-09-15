import {
  emptyArtist,
  emptySong,
  emptySongSettings,
  reduxRender,
  withToastify
} from '../../../../../../../test-utils.tsx'
import AddNewSongPart from './AddNewSongPart.tsx'
import { act, screen } from '@testing-library/react'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { userEvent } from '@testing-library/user-event'
import { Instrument } from '../../../../../../../types/models/Song.ts'
import { CreateSongPartRequest } from '../../parts/types/requests/SongPartRequests.ts'
import { BandMember } from '../../../../../../../types/models/Artist.ts'
import { setSong } from '../../../state/slice/songSlice.tsx'

describe('Add New Song Part', () => {
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
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    const [_, store] = reduxRender(<AddNewSongPart opened={true} onClose={() => {}} />, {
      song: { songId: '', isArtistBand: false, settings: emptySongSettings }
    })

    expect(screen.getByRole('button', { name: 'select-band-member' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'select-band-member' })).toBeDisabled()
    expect(screen.getByRole('button', { name: 'select-instrument' })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /name/i })).toBeInTheDocument()
    expect(await screen.findByRole('textbox', { name: /name/i })).toHaveFocus()
    expect(screen.getByRole('button', { name: /add/i })).toBeInTheDocument()

    expect(screen.getByRole('textbox', { name: /name/i })).not.toBeInvalid()

    const newSong = {
      ...emptySong,
      artist: { ...emptyArtist, isBand: true, bandMembers: bandMembers }
    }
    await act(() => store.dispatch(setSong(newSong)))
    expect(screen.getByRole('button', { name: 'select-band-member' })).not.toBeDisabled()
  })

  it('should have default options based on settings', async () => {
    const defaultInstrument = instruments[1]
    const defaultBandMember = bandMembers[1]

    const songSettings = { ...emptySongSettings, defaultBandMember, defaultInstrument }

    reduxRender(<AddNewSongPart opened={true} onClose={() => {}} />, {
      song: {
        songId: '',
        isArtistBand: true,
        artistBandMembers: bandMembers,
        settings: songSettings
      }
    })

    expect(screen.getByRole('button', { name: defaultBandMember.name })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: defaultInstrument.name })).toBeInTheDocument()
  })

  it('should send create request when name is typed', async () => {
    const user = userEvent.setup()

    const onClose = vitest.fn()
    const songId = 'some id'

    const newName = 'Part 1'

    let capturedRequest: CreateSongPartRequest
    server.use(
      http.post('/songs/parts', async (req) => {
        capturedRequest = (await req.request.json()) as CreateSongPartRequest
        return HttpResponse.json({ message: 'part added!' })
      })
    )

    reduxRender(withToastify(<AddNewSongPart opened={true} onClose={onClose} />), {
      song: {
        songId: songId,
        isArtistBand: true,
        artistBandMembers: bandMembers,
        settings: emptySongSettings
      }
    })

    await user.type(screen.getByRole('textbox', { name: /name/i }), newName)

    await user.click(screen.getByRole('button', { name: /add/i }))

    expect(capturedRequest).toStrictEqual({
      name: newName,
      songId: songId,
      sectionIds: []
    })
    expect(onClose).toHaveBeenCalledOnce()

    expect(screen.getByText(`${newName} added!`)).toBeInTheDocument()

    expect(screen.getByRole('textbox', { name: /name/i })).toHaveValue('')
  })

  it('should send create request when all fields are filled', async () => {
    const user = userEvent.setup()

    const onClose = vitest.fn()
    const songId = 'some id'

    const newName = 'Part 1'
    const newInstrument = instruments[0]
    const newBandMember = bandMembers[0]

    let capturedRequest: CreateSongPartRequest
    server.use(
      http.post('/songs/parts', async (req) => {
        capturedRequest = (await req.request.json()) as CreateSongPartRequest
        return HttpResponse.json({ message: 'part added!' })
      })
    )

    reduxRender(withToastify(<AddNewSongPart opened={true} onClose={onClose} />), {
      song: {
        songId: songId,
        isArtistBand: true,
        artistBandMembers: bandMembers,
        settings: emptySongSettings
      }
    })

    // fill fields
    await user.click(screen.getByRole('button', { name: 'select-band-member' }))
    await user.click(await screen.findByRole('option', { name: newBandMember.name }))

    await user.click(screen.getByRole('button', { name: 'select-instrument' }))
    await user.click(await screen.findByRole('option', { name: newInstrument.name }))

    await user.type(screen.getByRole('textbox', { name: /name/i }), newName)

    await user.click(screen.getByRole('button', { name: /add/i }))

    expect(capturedRequest).toStrictEqual({
      bandMemberId: newBandMember.id,
      instrumentId: newInstrument.id,
      name: newName,
      songId: songId,
      sectionIds: []
    })
    expect(onClose).toHaveBeenCalledOnce()

    expect(screen.getByText(`${newName} added!`)).toBeInTheDocument()

    // reset fields
    expect(screen.getByRole('button', { name: 'select-band-member' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'select-instrument' })).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /name/i })).toHaveValue('')
  })

  it('should send create request when there are default settings', async () => {
    const user = userEvent.setup()

    const onClose = vitest.fn()
    const songId = 'some id'
    const settings = {
      ...emptySongSettings,
      defaultBandMember: bandMembers[0],
      defaultInstrument: instruments[0]
    }

    const newName = 'Part 1'
    const newInstrument = instruments[1]

    let capturedRequest: CreateSongPartRequest
    server.use(
      http.post('/songs/parts', async (req) => {
        capturedRequest = (await req.request.json()) as CreateSongPartRequest
        return HttpResponse.json({ message: 'part added!' })
      })
    )

    reduxRender(withToastify(<AddNewSongPart opened={true} onClose={onClose} />), {
      song: {
        songId: songId,
        isArtistBand: true,
        artistBandMembers: bandMembers,
        settings: settings
      }
    })

    await user.click(screen.getByRole('button', { name: settings.defaultInstrument.name }))
    await user.clear(screen.getByRole('textbox', { name: /search/i }))
    await user.click(await screen.findByRole('option', { name: newInstrument.name }))

    await user.type(screen.getByRole('textbox', { name: /name/i }), newName)

    await user.click(screen.getByRole('button', { name: /add/i }))

    expect(capturedRequest).toStrictEqual({
      bandMemberId: settings.defaultBandMember.id,
      instrumentId: newInstrument.id,
      name: newName,
      songId: songId,
      sectionIds: []
    })
    expect(onClose).toHaveBeenCalledOnce()

    expect(screen.getByText(`${newName} added!`)).toBeInTheDocument()

    // reset fields
    expect(
      screen.getByRole('button', { name: settings.defaultBandMember.name })
    ).toBeInTheDocument()
    expect(
      screen.getByRole('button', { name: settings.defaultInstrument.name })
    ).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: /name/i })).toHaveValue('')
  })

  // Validation

  it('should display error when name was typed and then removed', async () => {
    const user = userEvent.setup()

    const newName = 'New Name'

    reduxRender(<AddNewSongPart opened={true} onClose={() => {}} />, {
      song: {
        songId: '',
        isArtistBand: true,
        artistBandMembers: bandMembers,
        settings: emptySongSettings
      }
    })

    await user.type(screen.getByRole('textbox', { name: /name/i }), newName)
    await user.clear(screen.getByRole('textbox', { name: /name/i }))

    expect(screen.getByRole('textbox', { name: /name/i })).toBeInvalid()
  })

  it('should not send create request when name is not typed', async () => {
    const user = userEvent.setup()

    let capturedRequest: CreateSongPartRequest
    server.use(
      http.post('/songs/parts', async (req) => {
        capturedRequest = (await req.request.json()) as CreateSongPartRequest
        return HttpResponse.json({ message: 'part added!' })
      })
    )

    reduxRender(<AddNewSongPart opened={true} onClose={() => {}} />, {
      song: {
        songId: '',
        isArtistBand: true,
        artistBandMembers: bandMembers,
        settings: emptySongSettings
      }
    })

    await user.click(screen.getByRole('button', { name: /add/i }))

    expect(screen.getByRole('textbox', { name: /name/i })).toBeInvalid()

    expect(capturedRequest).toBeUndefined()
  })

  it('should refresh errors when reopened', async () => {
    const user = userEvent.setup()

    const uut = (opened = true) => <AddNewSongPart opened={opened} onClose={() => {}} />

    const [{ rerender }] = reduxRender(uut(), {
      song: {
        songId: '',
        isArtistBand: true,
        artistBandMembers: bandMembers,
        settings: emptySongSettings
      }
    })

    await user.click(screen.getByRole('button', { name: /add/i }))

    expect(screen.getByRole('textbox', { name: /name/i })).toBeInvalid()

    rerender(uut(false))
    rerender(uut())

    expect(screen.getByRole('textbox', { name: /name/i })).not.toBeInvalid()
  })
})
