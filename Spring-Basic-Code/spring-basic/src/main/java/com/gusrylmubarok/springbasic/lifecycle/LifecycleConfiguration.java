package com.gusrylmubarok.springbasic.lifecycle;

import com.gusrylmubarok.springbasic.data.Connection;
import com.gusrylmubarok.springbasic.data.cyclic.Server;
import org.springframework.context.annotation.Bean;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

public class LifecycleConfiguration {
    @Bean
    public Connection connection(){
        return new Connection();
    }

    // @Bean(initMethod = "start", destroyMethod = "stop")
    @Bean
    public Server server(){
        return new Server();
    }

}
