import { mantineRender } from '../test-utils.tsx'
import Unauthorized from "./Unauthorized.tsx";
import {screen} from "@testing-library/react";

describe('Unauthorized', () => {
  it('should render', () => {
    mantineRender(<Unauthorized />)

    expect(screen.getByRole('heading', { name: /unauthorized/i })).toBeInTheDocument()
  })
})
