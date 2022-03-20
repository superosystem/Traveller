package com.gusrylmubarok.springbasic;

import com.gusrylmubarok.springbasic.data.Connection;
import com.gusrylmubarok.springbasic.data.cyclic.Server;
import com.gusrylmubarok.springbasic.lifecycle.LifecycleConfiguration;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.context.ConfigurableApplicationContext;
import org.springframework.context.annotation.AnnotationConfigApplicationContext;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

public class LifeCycleTest {

    private ConfigurableApplicationContext applicationContext;

    @BeforeEach
    void setUp() {
        applicationContext = new AnnotationConfigApplicationContext(LifecycleConfiguration.class);
        applicationContext.registerShutdownHook();
    }

    @AfterEach
    void tearDown() {
        // applicationContext.close();
    }

    @Test
    void testConnection() {

        Connection connection = applicationContext.getBean(Connection.class);

    }

    @Test
    void testServer() {
        Server server = applicationContext.getBean(Server.class);
    }
}
