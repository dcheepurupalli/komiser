import { ReactNode, useContext, useState } from 'react';


type InfrastructureLayoutProps = {
    children: ReactNode;
}

function InfrastructureLayout({ children }: InfrastructureLayoutProps) {
    return (
        <>
        <main>
            {children}
        </main>
        </>
    )
}

export default InfrastructureLayout;