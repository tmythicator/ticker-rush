import styles from './Legal.module.css';

export const AgbPage = () => {
  return (
    <div className={styles.container}>
      <h1 className={styles.title}>General Terms and Conditions (AGB)</h1>

      <section className={styles.section}>
        <h2 className={styles.sectionHeader}>1. Scope of Application</h2>
        <p className={styles.text}>
          These General Terms and Conditions (GTC / AGB) apply to all users of the Ticker Rush
          platform. By registering an account, you agree to these terms.
        </p>
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionHeader}>2. Services Provided</h2>
        <p className={styles.text}>
          Ticker Rush is a private project designed for demonstration and educational purposes. It
          provides a simulated trading environment utilizing mock currency. There is no real-world
          financial value, and no real money is deposited, traded, or withdrawn.
        </p>
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionHeader}>3. User Obligations</h2>
        <p className={styles.text}>
          Users must provide accurate information during registration. You are responsible for
          keeping your login credentials confidential. The platform must not be used for any
          unlawful activities or automated scraping without permission.
        </p>
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionHeader}>4. Liability</h2>
        <p className={styles.text}>
          As a free, educational simulation, we assume no liability for the accuracy, completeness,
          or timeliness of the market data provided. The platform is offered "as is" without any
          warranty of continuous availability. We are not liable for any damages resulting from the
          use or inability to use the service.
        </p>
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionHeader}>5. Termination</h2>
        <p className={styles.text}>
          We reserve the right to suspend or terminate user accounts that violate these terms or
          engage in disruptive behavior. Users may stop using the service at any time.
        </p>
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionHeader}>6. Changes to the AGB</h2>
        <p className={styles.text}>
          We reserve the right to amend these terms. Significant changes will be communicated to
          users, and continued use of the platform following modifications implies acceptance of the
          new AGB.
        </p>
      </section>
    </div>
  );
};
