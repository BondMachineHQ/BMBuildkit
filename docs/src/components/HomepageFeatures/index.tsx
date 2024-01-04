import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';

type FeatureItem = {
    title: string;
    Svg: React.ComponentType<React.ComponentProps<'svg'>>;
    description: JSX.Element;
};

const FeatureList: FeatureItem[] = [
    {
        title: 'Docker-like simple!',
        Svg: require('@site/static/img/home-2.svg').default,
        description: (
            <>
                Designed for emulating the experience of a container build via Docker.
                You define a BMFile and an FPGA board, the rest is on us.
            </>
        ),
    },
    {
        title: 'Program FPGA on the cloud... if you want to',
        Svg: require('@site/static/img/home-1.svg').default,
        description: (
            <>
                Get your firmware from any container registry and program your FPGA with a single command.
                And yes, CI/CD is never been easier.
            </>
        ),
    },
    {
        title: 'Prototype as you want!',
        Svg: require('@site/static/img/home-3.svg').default,
        description: (
            <>
                Your firmwares will be tagged like a usual container image, thus keeping track of all your experiments.
                Comparing the performances will be a walk in the park.
            </>
        ),
    },
];

function Feature({title, Svg, description}: FeatureItem) {
    return (
        <div className={clsx('col col--4')}>
            <div className="text--center">
                <Svg className={styles.featureSvg} role="img"/>
            </div>
            <div className="text--center padding-horiz--md">
                <Heading as="h3">{title}</Heading>
                <p>{description}</p>
            </div>
        </div>
    );
}

export default function HomepageFeatures(): JSX.Element {
    return (
        <section className={styles.features}>
            <div className="container">
                <div className="row">
                    {FeatureList.map((props, idx) => (
                        <Feature key={idx} {...props} />
                    ))}
                </div>
            </div>
        </section>
    );
}
