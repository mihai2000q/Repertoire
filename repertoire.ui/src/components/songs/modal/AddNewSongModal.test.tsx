import { reduxRender } from '../../../test-utils.tsx'
import AddNewSongModal from './AddNewSongModal.tsx'
import { vi } from 'vitest'
import { http, HttpResponse } from 'msw'
import { setupServer } from 'msw/node'
import { screen } from '@testing-library/react'
import { userEvent } from '@testing-library/user-event'
import { CreateSongRequest } from '../../../types/requests/SongRequests.ts'

describe.skip('Add New Song Modal', () => {
  let capturedCreateSongRequest: CreateSongRequest | undefined

  const handlers = [
    http.post('/songs', async (req) => {
      capturedCreateSongRequest = (await req.request.json()) as CreateSongRequest
      return HttpResponse.json({ id: 'some id' })
    })
  ]

  const server = setupServer(...handlers)

  beforeAll(() => server.listen())

  afterEach(() => {
    capturedCreateSongRequest = undefined
    server.resetHandlers()
  })

  afterAll(() => server.close())

  it('should render and display form', ({ expect }) => {
    reduxRender(<AddNewSongModal opened={true} onClose={vi.fn()} />)

    expect(screen.getByRole('heading', { name: /add new song/i })).toBeInTheDocument()

    expect(screen.getByRole('textbox', { name: /title/i })).toHaveTextContent('')
    expect(screen.getByRole('button', { name: /add-image-button/i })).toBeInTheDocument()

    expect(screen.getByRole('button', { name: /add song/i })).toBeInTheDocument()
  })

  it('should display error if the title is invalid', async ({ expect }) => {
    const user = userEvent.setup()
    const error = 'Title cannot be blank'

    reduxRender(<AddNewSongModal opened={true} onClose={vi.fn()} />)

    await user.click(screen.getByRole('button', { name: /add song/i }))
    expect(screen.getByText(error)).toBeInTheDocument()

    await user.type(screen.getByRole('textbox', { name: /title/i }), '  ')
    await user.click(screen.getByRole('button', { name: /add song/i }))
    expect(screen.getByText(error)).toBeInTheDocument()
  })

  it('should send POST request when valid', async ({ expect }) => {
    const onClose = vi.fn()
    const user = userEvent.setup()
    const title = 'New Title'

    reduxRender(<AddNewSongModal opened={true} onClose={onClose} />)

    await user.type(screen.getByRole('textbox', { name: /title/i }), title)
    await user.click(screen.getByRole('button', { name: /add song/i }))

    expect(capturedCreateSongRequest).toStrictEqual({
      title: title
    })
    expect(onClose).toHaveBeenCalledOnce()
    expect(screen.getByRole('textbox', { name: /title/i })).toHaveTextContent('')
  })
})
