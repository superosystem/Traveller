package com.gusrylmubarok.hibernate;

import com.gusrylmubarok.hibernate.config.HibernateConfiguration;
import junit.framework.TestCase;
import lombok.extern.slf4j.Slf4j;
import org.hibernate.Session;
import org.junit.Test;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

@Slf4j
public class HibernateConfigurationTest extends TestCase {
    private Session session;


    @Override
    protected void setUp() throws Exception {
        log.info("START HIBERNATE SESSION");
        this.session = HibernateConfiguration.getSession();
    }

    @Test
    public void testOpenConnection() {
        this.session.beginTransaction();
    }

    @Override
    protected void tearDown() throws Exception {
        log.info("DESTROY HIBERNATE SESSION");
        this.session.close();
    }
}
