import { SongArrangement } from '../../../types/models/Song.ts'
import { emptySongArrangement, reduxRender, withToastify } from '../../../test-utils.tsx'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { userEvent } from '@testing-library/user-event/dist/cjs/index.js'
import { expect } from 'vitest'
import { screen } from '@testing-library/react'
import CustomRehearsalButton from './CustomRehearsalButton.tsx'
import { AddCustomSongRehearsalRequest } from '../../../types/requests/SongRequests.ts'

describe('Custom Rehearsal Button', () => {
  const arrangements: SongArrangement[] = [
    {
      ...emptySongArrangement,
      id: '1',
      name: 'something'
    },
    {
      ...emptySongArrangement,
      id: '2',
      name: 'something 2'
    }
  ]

  const handlers = [
    http.get(`/songs/arrangements`, () => {
      return HttpResponse.json(arrangements)
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => server.resetHandlers())

  afterAll(() => server.close())

  it('should render', async () => {
    const user = userEvent.setup()

    reduxRender(<CustomRehearsalButton songId={''} sectionsCount={1} />)

    expect(screen.getByRole('button', { name: 'add-custom-rehearsal' })).toBeInTheDocument()
    expect(await screen.findByRole('button', { name: 'add-custom-rehearsal' })).not.toBeDisabled()
    await user.click(screen.getByRole('button', { name: 'add-custom-rehearsal' }))

    for (const arrangement of arrangements) {
      expect(await screen.findByRole('menuitem', { name: arrangement.name })).toBeInTheDocument()
    }

    await user.click(screen.getByRole('menuitem', { name: arrangements[0].name }))

    expect(await screen.findByRole('button', { name: 'cancel' })).toBeInTheDocument()
    expect(await screen.findByRole('button', { name: 'confirm' })).toBeInTheDocument()
  })

  it("should be disabled when sections' count is 0", () => {
    reduxRender(<CustomRehearsalButton songId={''} sectionsCount={0} />)

    expect(screen.getByRole('button', { name: 'add-custom-rehearsal' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'add-custom-rehearsal' })).toBeDisabled()
  })

  it('should be disabled when there are no arrangements', async () => {
    server.use(
      http.get(`/songs/arrangements`, () => {
        return HttpResponse.json([])
      })
    )

    reduxRender(<CustomRehearsalButton songId={''} sectionsCount={1} />)

    expect(await screen.findByRole('button', { name: 'add-custom-rehearsal' })).toBeDisabled()
  })

  it('should send custom rehearsal request when confirming on an arrangement', async () => {
    const user = userEvent.setup()

    let capturedRequest: AddCustomSongRehearsalRequest
    server.use(
      http.post('/songs/custom-rehearsal', async (req) => {
        capturedRequest = (await req.request.json()) as AddCustomSongRehearsalRequest
        return HttpResponse.json({ message: 'it worked' })
      })
    )

    const arrangement = arrangements[0]
    const songId = 'some-id'

    reduxRender(withToastify(<CustomRehearsalButton songId={songId} sectionsCount={1} />))

    await user.click(screen.getByRole('button', { name: 'add-custom-rehearsal' }))

    await user.click(await screen.findByRole('menuitem', { name: arrangement.name }))
    await user.click(await screen.findByRole('button', { name: 'confirm' }))

    expect(await screen.findByText(/custom rehearsal added/i)).toBeInTheDocument()
    expect(capturedRequest).toStrictEqual({ id: songId, arrangementId: arrangement.id })
  })
})
